package main

// All the bytes supplied here will be either +, - or 0...9
func bytesToInt(byteArray []byte) int {
	var result int
	negative := false

	for _, b := range byteArray {
		if b == 45 { // 45 = "-" signal, as per ASCII conventions
			negative = true
			continue
		}
		if b == 46 { // skip the decimal point
			continue
		}
		// For each new number, move the old number one digit to left; byte shifting
		result = result*10 + int(b-48) // 48 = '0', 49 = '1', ...
	}

	// handle negative numbers
	if negative {
		return -result
	}
	return result
}
