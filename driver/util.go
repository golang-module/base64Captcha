package driver

func stringToFakeByte(content string) []byte {
	digits := make([]byte, len(content))
	for idx, cc := range content {
		digits[idx] = byte(cc - '0')
	}
	return digits
}

// parseDigitsToString parse randomDigits to normal string
func parseDigitsToString(bytes []byte) string {
	stringB := make([]byte, len(bytes))
	for idx, by := range bytes {
		stringB[idx] = by + '0'
	}
	return string(stringB)
}
