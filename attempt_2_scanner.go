package main

import (
	"bufio"
	"os"
)

// Just speed reading with a basic scanner and nothing else
func attemptTwo() {
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
