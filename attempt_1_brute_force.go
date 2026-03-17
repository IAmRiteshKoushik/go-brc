package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func attemptOne() {
	data := make(map[string]*StationData)

	file, err := os.Open("measurements.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		name := parts[0]
		tempStr := strings.Trim(parts[1], "\n") // Removing the newline character

		temperature, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			panic(err)
		}

		station, ok := data[name]
		// If the station does not exist, create a new station
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
	printResult(data)
}
