package main

import (
	"bufio"
	"bytes"
	"hash/fnv"
	"sync"
)

func consumerV5(input chan []byte, output chan map[uint64]*StationDataV2, wg *sync.WaitGroup) {
	defer wg.Done()

	// This time instead of using a string as a key, I am going to use a uint64
	// This is because, I can hash the bytes and generate this pretty easily.
	// So, the memory allocation of the has
	data := make(map[uint64]*StationDataV2)

	hash := fnv.New64a() // a hash function that returns uint64
	nameBuf := make([]byte, 100)
	tempBuf := make([]byte, 50)

	for reading := range input {
		scanner := bufio.NewScanner(bytes.NewReader(reading))
		for scanner.Scan() {
			// Operate directly on bytes instead of converting them into strings
			line := scanner.Bytes()
			nameSize, tempSize := ParseLine(line, nameBuf, tempBuf)
			name := nameBuf[:nameSize] // No additional allocs, only ptr

			// name := string(nameBuf[:nameSize]) // PERFORMANCE KILLER

			// PERFORMANCE KILLER
			// temperature, err := strconv.ParseFloat(string(tempBuf[:tempSize]), 64)

			temperature := bytesToInt(tempBuf[:tempSize])

			hash.Reset()
			hash.Write(name)
			id := hash.Sum64() // generates the uint64 key

			station, ok := data[id]
			if !ok {
				// Name conversion is only done when a new entry is created into the hashmap
				// otherwise we never see that byte to string conversion.
				data[id] = &StationDataV2{string(name), temperature, temperature, temperature, 1}
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
