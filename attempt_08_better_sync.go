package main

import (
	"io"
	"os"
)

func attemptEight(BUFFER_SIZE int) {
	channel := make(chan []byte, 10)
	go consumer(channel)

	file, err := os.Open("measurements.txt")
	if err != nil {
	}

	buffer := make([]byte, BUFFER_SIZE)

	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}

		// Instead of sending the buffer directly, copy a part of the buffer and
		// send it to the consumer. Reasons for doing this is explained in the notes
		data := make([]byte, BUFFER_SIZE)
		copy(data, buffer[:BUFFER_SIZE])
		channel <- data
	}
}
