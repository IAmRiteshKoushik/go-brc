package main

import (
	"fmt"
	"io"
	"os"
)

func attemptSeven(BUFFER_SIZE int) {
	channel := make(chan []byte, 10)
	go consumer(channel)

	file, err := os.Open("measurements.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, BUFFER_SIZE)
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			return
		}

		if err != nil {
			panic(err)
		}

		channel <- buffer
	}
}
