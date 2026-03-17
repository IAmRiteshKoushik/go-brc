package main

import (
	"bufio"
	"io"
	"os"
)

// reading with Reader instead of scanner
func attemptFive() {
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		_, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}
