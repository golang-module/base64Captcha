package driver

import (
	"crypto/rand"
	"image/color"
	"io"
	mathRand "math/rand"
	"strings"
)

const (
	StringLength = 20
	TxtNumbers   = "012346789"
	TxtAlphabet  = "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz"
)

func RandomInt(n int) int {
	if n > 0 {
		return mathRand.Intn(n)
	}
	return 0
}

func RandomString() string {
	var bytes = []byte(TxtNumbers + TxtAlphabet)
	b := randomBytesMod(StringLength, byte(len(bytes)))
	for i, c := range b {
		b[i] = bytes[c]
	}
	return string(b)
}

// RandomText creates random text of given size.
func RandomText(size int, sourceChars string) string {
	if sourceChars == "" || size == 0 {
		return ""
	}

	if size >= len(sourceChars) {
		sourceChars = strings.Repeat(sourceChars, size)
	}

	sourceRunes := []rune(sourceChars)
	sourceLength := len(sourceRunes)

	text := make([]rune, size)
	for i := range text {
		text[i] = sourceRunes[RandomInt(sourceLength)]
	}
	return string(text)
}

// RandomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-9. The slice can be used as a captcha
// solution.
func RandomDigits(length int) []byte {
	return randomBytesMod(length, 10)
}

// RandomBytes returns a byte slice of the given length read from CSPRNG.
func RandomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

func RandomColor() color.RGBA {
	red := RandomInt(55) + 200
	green := RandomInt(55) + 200
	blue := RandomInt(55) + 200
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxRB := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := RandomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxRB {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}
}
