package Utils

import (
	"strconv"
)

func isRuneADigit(r uint8) bool {
	// 47 - '0'
	// 57 - '9'
	return r >= 47 && r <= 57
}

func IncrementName(name string) string {
	// First, we need to read all the digits from the end of the string
	lastNum := ""
	for i := len(name) - 1; i >= 0; i-- {
		// If it's a digit, prepend the lastNum with it
		if isRuneADigit(name[i]) {
			lastNum = string(name[i]) + lastNum
		} else {
			// Else break the loop
			break
		}
	}
	var n int
	if len(lastNum) > 0 {
		n, _ = strconv.Atoi(lastNum)
	} else {
		n = 0
	}
	n += 1
	return name[:len(name)-len(lastNum)] + strconv.Itoa(n)
}
