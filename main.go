package main

import (
	"fmt"
	"os"
	"sort"
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
	start := time.Now()
	switch flag {
	case "one":
		attemptOne() // read and compute
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "two":
		attemptTwo() // read-only
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "three":
		for range 8 {
			attemptThree(1024 * 1024) // read-only
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "four":
		attemptFour() // read-only, using a reader.ReadBytes instead of a scanner
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "five":
		attemptFive() // read-only, using a reader.ReadLine instead of a scanner
		fmt.Printf("%0.6f\n", time.Since(start).Seconds())
	case "six":
		for range 4 {
			attemptSix() // read-only, using a file.Read
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	case "seven":
		for range 4 {
			attemptSeven(512 * 512)
			fmt.Printf("%0.6f\n", time.Since(start).Seconds())
		}
	}
}

func printResult(data map[string]*StationData) {
	result := make(map[string]*StationData, len(data))
	keys := make([]string, 0, len(data))
	for _, v := range data {
		keys = append(keys, v.Name)
		result[v.Name] = v
	}
	// We are sorting things alphabetically and then printing, this adds a bit of
	// additional overhead especially when processing 40k unique names
	sort.Strings(keys)

	print("{\n")
	for _, k := range keys {
		v := result[k]

		// Printing <station-name>=min/avg/max
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", k, v.Min, v.Sum/float64(v.Count), v.Max)
	}
	print("}\n")
}
