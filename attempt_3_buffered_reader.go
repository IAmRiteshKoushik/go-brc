package main

import (
	"bufio"
	"os"
)

// Just speed reading with a basic scanner and nothing else
func attemptThree(BUFFER_SIZE int) {
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, BUFFER_SIZE), BUFFER_SIZE)
	for scanner.Scan() {
		scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
