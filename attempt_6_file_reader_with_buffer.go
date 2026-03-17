package main

import (
	"io"
	"os"
)

// reading with Reader instead of scanner
func attemptSix() {
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024)
	for {
		_, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}
