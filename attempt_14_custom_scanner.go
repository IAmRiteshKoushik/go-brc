package main

import (
	"io"
	"os"
	"sync"
)

// The only difference is that we are now using a custom scanner
func attempFourteen(workers, chanBufSize int) {
	inputChannels := make([]chan []byte, workers)
	outputChannels := make([]chan map[uint64]*StationDataV2, workers)

	var wg sync.WaitGroup
	wg.Add(workers)

	// Spawn workers
	for i := range workers {
		input := make(chan []byte, chanBufSize)
		output := make(chan map[uint64]*StationDataV2, 1)

		go consumerV6(input, output, &wg)

		inputChannels[i] = input
		outputChannels[i] = output
	}

	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// readBuffer best size is already determined to be 512 * 512
	readBuffer := make([]byte, 512*512)
	leftoverBuffer := make([]byte, 1024) // sufficient size of pending buffers
	leftoverSize := 0
	currentWorker := 0

	for {
		// You have the number of bytes read, despite the buffer size. You will not
		// always use the entire buffer
		n, err := file.Read(readBuffer)
		if err == io.EOF {
			// Send
			if leftoverSize > 0 {
				data := make([]byte, leftoverSize)
				copy(data, leftoverBuffer[:leftoverSize])
				inputChannels[currentWorker] <- data
			}
			break
		}
		if err != nil {
			panic(err)
		}

		m := 0
		// Finding the last newline character and capturing its index
		for i := n - 1; i >= 0; i-- {
			if readBuffer[i] == 10 {
				m = i
				break
			}
		}

		data := make([]byte, m+leftoverSize)
		copy(data, leftoverBuffer[:leftoverSize])
		copy(data[leftoverSize:], readBuffer[:m])
		copy(leftoverBuffer, readBuffer[m+1:n]) // prep for next iteration
		leftoverSize = n - m - 1

		inputChannels[currentWorker] <- data

		currentWorker++
		if currentWorker >= workers {
			currentWorker = 0
		}
	}

	// close the input channels, make the workers leave their processing loop
	// This is fine to do as all reading is done.
	for i := range workers {
		close(inputChannels[i])
	}

	// Wait for all workers to finish processing if they are still in transit
	// Once that is done, then close the output channels too
	wg.Wait()
	// for i := range workers {
	// 	close(outputChannels[i])
	// }

	// Now we can start aggregating
	data := make(map[uint64]*StationDataV2)
	for i := range workers {
		// Extract stuff out of channels synchronously
		for stationHash, stationData := range <-outputChannels[i] {
			// Check for existance of record
			if _, ok := data[stationHash]; !ok {
				data[stationHash] = stationData
			} else {
				// If record already exists, then handle the comparison and increments
				if stationData.Min < data[stationHash].Min {
					data[stationHash].Min = stationData.Min
				}
				if stationData.Max > data[stationHash].Max {
					data[stationHash].Max = stationData.Max
				}
				data[stationHash].Sum += stationData.Sum
				data[stationHash].Count += stationData.Count
			}
		}
	}
}
