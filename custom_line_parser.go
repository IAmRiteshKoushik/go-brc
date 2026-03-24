package main

func ParseLine(line, nameBuffer, tempBuffer []byte) (nameSize, tempSize int) {
	// 59 is the ASCII character for ";"
	// 10 is the ASCII character for "\n"
	total := len(line)
	i := 0

	// Determine the name size and  populate the corresponding byte slices
	j := 0
	for line[i] != 59 {
		nameBuffer[j] = line[i]
		i++
		j++
	}

	i++ // skip the semi-colon

	// Determine temperature size and populate the corresponding byte slices
	k := 0
	for i < total && line[i] != 10 {
		tempBuffer[k] = line[i]
		i++
		k++
	}

	return j, k // sizes
}
