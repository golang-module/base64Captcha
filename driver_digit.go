package captcha

import "math/rand"

// DriverDigit config for captcha-engine-digit.
type DriverDigit struct {
	// Height png height in pixel.
	Height int
	// Width Captcha png width in pixel.
	Width int
	// DefaultLen Default number of digits in captcha solution.
	Length int
	// MaxSkew max absolute skew factor of a single digit.
	MaxSkew float64
	// DotCount Number of background circles.
	DotCount int
}

// NewDriverDigit creates a driver of digit
func NewDriverDigit(height int, width int, length int, maxSkew float64, dotCount int) *DriverDigit {
	return &DriverDigit{Height: height, Width: width, Length: length, MaxSkew: maxSkew, DotCount: dotCount}
}

// DefaultDriverDigit is a default driver of digit
var DefaultDriverDigit = NewDriverDigit(80, 240, 5, 0.7, 80)

// GenerateIdQuestionAnswer creates captcha content and answer
func (d *DriverDigit) GenerateIdQuestionAnswer() (id, q, a string) {
	id = RandomId()
	digits := randomDigits(d.Length)
	a = parseDigitsToString(digits)
	return id, a, a
}

// GenerateSpecificIdQuestionAnswer creates captcha content and answer
func (d *DriverDigit) GenerateSpecificIdQuestionAnswer(mId string) (id, q, a string) {
	id = mId
	digits := randomDigits(d.Length)
	a = parseDigitsToString(digits)
	return id, a, a
}

// DrawCaptcha creates digit captcha item
func (d *DriverDigit) DrawCaptcha(content string) (item Item, err error) {
	// Initialize PRNG.
	itemDigit := NewItemDigit(d.Width, d.Height, d.DotCount, d.MaxSkew)
	// parse digits to string
	digits := stringToFakeByte(content)

	itemDigit.calculateSizes(d.Width, d.Height, len(digits))
	// Randomly position captcha inside the image.
	maxX := d.Width - (itemDigit.width+itemDigit.dotSize)*len(digits) - itemDigit.dotSize
	maxY := d.Height - itemDigit.height - itemDigit.dotSize*2
	var border int
	if d.Width > d.Height {
		border = d.Height / 5
	} else {
		border = d.Width / 5
	}
	x := randIntn(maxX-border*2) + border
	y := randIntn(maxY-border*2) + border
	// Draw digits.
	for _, n := range digits {
		itemDigit.drawDigit(digitFontData[n], x, y)
		x += itemDigit.width + itemDigit.dotSize
	}
	// Draw strike-through line.
	itemDigit.strikeThrough()
	// Apply wave distortion.
	itemDigit.distort(rand.Float64()*(10-5)+5, rand.Float64()*(200-100)+100)
	// Fill image with random circles.
	itemDigit.fillWithCircles(d.DotCount, itemDigit.dotSize)
	return itemDigit, nil
}
