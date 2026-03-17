package main

import (
	"io"
	"os"
)

const bufferSize = 512 * 512

func attemptNine() {
	channel := make(chan []byte, 10)
	go consumer(channel)

	file, err := os.Open("measurements.txt")
	if err != nil {
	}

	// Out of loop's scope
	readBuffer := make([]byte, bufferSize)
	leftoverBuffer := make([]byte, 1024) // a little bit extra
	leftoverSize := 0

	// Read the file in chunks and send it via a channel to consumer
	for {
		n, err := file.Read(readBuffer)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}

		// Find the last newline character "\n" (byte = 10). Loop from end to
		// first and then find it
		m := 0
		for i := bufferSize - 1; i >= 0; i-- {
			if readBuffer[i] == 10 {
				m = i
				break
			}
		}

		// Byte buffer is created base current valid size and previous leftover
		data := make([]byte, m+leftoverSize)

		// data buffer is populated based on previous leftover and the next complete
		// byte set
		copy(data, leftoverBuffer[:leftoverSize])
		copy(data[leftoverSize:], readBuffer[:m])

		// Leftover buffer for the next iteration with incomplete byte set
		copy(leftoverBuffer, readBuffer[m+1:bufferSize])
		leftoverSize = n - m - 1

		channel <- data
	}
}
