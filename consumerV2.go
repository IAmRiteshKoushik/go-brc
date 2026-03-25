package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"sync"
)

// A more complicated consumer, this one reads lines; captures data and sends
// it back to the main thread over a channel of map[string]*StationData. Using
// the pointer to reduce mallocs
func consumerV2(input chan []byte, output chan map[string]*StationData, wg *sync.WaitGroup) {
	defer wg.Done()
	data := make(map[string]*StationData)

	separator := []byte{';'}

	for reading := range input {
		scanner := bufio.NewScanner(bytes.NewReader(reading))
		for scanner.Scan() {
			// Operate directly on bytes instead of converting them into strings
			line := scanner.Bytes()
			parts := bytes.Split(line, separator)

			if len(parts) != 2 {
				fmt.Println("Invalid line: ", string(line))
				continue
			}

			name := string(parts[0])
			temperature, err := strconv.ParseFloat(string(parts[1]), 64)
			if err != nil {
				fmt.Println(err)
				return
			}

			station, ok := data[name]
			if !ok {
				data[name] = &StationData{name, temperature, temperature, temperature, 1}
			} else {
				if temperature < station.Min {
					station.Min = temperature
				}
				if temperature > station.Max {
					station.Max = temperature
				}
				station.Sum += temperature
				station.Count++
			}
		}
	}

	output <- data
}
