package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func run_bytes() {
	/*
	 * Reading through the entire file using an internal buffer
	 * */
	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// String conversion is slower due to additional overhead as it involved
	// allocation of memory
	for scanner.Scan() {
		// Use an internal buffer returning the same object so there is no additional
		// allocation
		scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func run_strings() {
	/*
	 * Reading through the entire file line by line in the form of strings
	 * */

	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Text()
	}
}

// Driver program
func main() {
	started := time.Now()
	run_bytes()
	fmt.Printf("Time taken to read 13GB file in byte-buffer: %.6f\n", time.Since(started).Seconds())
	started = time.Now()
	run_strings()
	fmt.Printf("Time taken to read 13GB file in strings: %.6f", time.Since(started).Seconds())
}
