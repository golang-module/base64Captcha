package digit

import (
	"github.com/golang-module/base64Captcha/driver"
	mathRand "math/rand"
)

// parseDigitsToString parse randomDigits to normal string
func parseDigitsToString(bytes []byte) string {
	stringB := make([]byte, len(bytes))
	for idx, by := range bytes {
		stringB[idx] = by + '0'
	}
	return string(stringB)
}
func stringToFakeByte(content string) []byte {
	digits := make([]byte, len(content))
	for idx, cc := range content {
		digits[idx] = byte(cc - '0')
	}
	return digits
}

func randIntRange(from, to int) int {
	if to-from <= 0 {
		return from
	}
	return driver.RandomInt(to-from) + from
}

func randFloat64Range(from, to float64) float64 {
	return mathRand.Float64()*(to-from) + from
}
