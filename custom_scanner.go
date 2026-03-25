package main

// This is derieved from the parseLine function
func nextLine(readingIdx int, reading, nameBuffer, tempBuffer []byte) (nexReadingIdx, nameSize, tempSize int) {
	i := readingIdx
	j := 0
	for reading[i] != 59 { // ;
		nameBuffer[j] = reading[i]
		i++
		j++
	}

	i++ // skip ;

	k := 0
	for i < len(reading) && reading[i] != 10 { // \n
		if reading[i] == 46 { // skipping decimal point
			i++
			continue
		}
		tempBuffer[k] = reading[i]
		i++
		k++
	}

	readingIdx = i + 1
	return readingIdx, j, k
}
