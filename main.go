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
	case "one":
		start := time.Now()
		attemptOne() // read and compute
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "two":
		start := time.Now()
		attemptTwo() // read-only
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "three":
		for range 8 {
			start := time.Now()
			attemptThree(1024 * 1024) // read-only
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "four":
		start := time.Now()
		attemptFour() // read-only, using a reader.ReadBytes instead of a scanner
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "five":
		start := time.Now()
		attemptFive() // read-only, using a reader.ReadLine instead of a scanner
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "six":
		for range 4 {
			start := time.Now()
			attemptSix() // read-only, using a file.Read
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "seven":
		for range 4 {
			start := time.Now()
			attemptSeven(512 * 512)
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "eight":
		for range 4 {
			start := time.Now()
			attemptEight(512 * 512)
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "nine":
		for range 4 {
			start := time.Now()
			attemptNine()
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	}
}
