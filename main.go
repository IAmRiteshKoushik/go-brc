package main

import (
	"fmt"
	"os"
	"time"
)

type StationData struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func main() {
	flag := os.Args[1]
	switch flag {
	case "1":
		start := time.Now()
		attemptOne() // read and compute
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "2":
		start := time.Now()
		attemptTwo() // read-only
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "3":
		for range 8 {
			start := time.Now()
			attemptThree(1024 * 1024) // read-only
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "4":
		start := time.Now()
		attemptFour() // read-only, using a reader.ReadBytes instead of a scanner
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "5":
		start := time.Now()
		attemptFive() // read-only, using a reader.ReadLine instead of a scanner
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "6":
		for range 4 {
			start := time.Now()
			attemptSix() // read-only, using a file.Read
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "7":
		for range 4 {
			start := time.Now()
			attemptSeven(512 * 512)
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "8":
		for range 4 {
			start := time.Now()
			attemptEight(512 * 512)
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "9":
		for range 4 {
			start := time.Now()
			attemptNine()
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "10":
		start := time.Now()
		// Workers, Buffer
		workerCount := 5
		bufferSize := 10
		attemptTen(workerCount, bufferSize)
		fmt.Printf("\t%0.2f", time.Since(start).Seconds())
	case "11":
	case "12":
	}
}
