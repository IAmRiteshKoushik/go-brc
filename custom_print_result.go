package main

import (
	"fmt"
	"sort"
)

func printResult(data map[string]*StationDataV2) {
	result := make(map[string]*StationDataV2, len(data))
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
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", k, v.Min/10.0, (v.Sum/v.Count)/10.0, v.Max/10.0)
	}
	print("}\n")
}
