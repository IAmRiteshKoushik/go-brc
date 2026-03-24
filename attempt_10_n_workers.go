package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func attemptTen(workers, chanBufSize int) {
	inputChannels := make([]chan []byte, workers)
	outputChannels := make([]chan map[string]*StationData, workers)

	var wg sync.WaitGroup
	wg.Add(workers)

	// Spawn workers
	for i := range workers {
		input := make(chan []byte, chanBufSize)
		output := make(chan map[string]*StationData, 1)

		go consumerV2(input, output, &wg)

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
		fmt.Printf("Work done by: %d\n", currentWorker)

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
	data := make(map[string]*StationData)
	for i := range workers {
		// Extract stuff out of channels synchronously
		for station, stationData := range <-outputChannels[i] {
			// Check for existance of record
			if _, ok := data[station]; !ok {
				data[station] = stationData
			} else {
				// If record already exists, then handle the comparison and increments
				if stationData.Min < data[station].Min {
					data[station].Min = stationData.Min
				}
				if stationData.Max > data[station].Max {
					data[station].Max = stationData.Max
				}
				data[station].Sum += stationData.Sum
				data[station].Count += stationData.Count
			}
		}
	}
}
