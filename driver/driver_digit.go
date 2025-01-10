package driver

import (
	"math/rand/v2"
)

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
func NewDriverDigit(d DriverDigit) *DriverDigit {
	return mergeDriverDigit(d)
}

// mergeDriverAudio merges default driver with given digit driver
func mergeDriverDigit(d DriverDigit) *DriverDigit {
	if d.Height == 0 {
		d.Height = DefaultDriverDigit.Height
	}
	if d.Width == 0 {
		d.Width = DefaultDriverDigit.Width
	}
	if d.Length == 0 {
		d.Length = DefaultDriverDigit.Length
	}
	if d.MaxSkew == 0 {
		d.MaxSkew = DefaultDriverDigit.MaxSkew
	}
	if d.DotCount == 0 {
		d.DotCount = DefaultDriverDigit.DotCount
	}
	return &d
}

// DrawCaptcha creates digit captcha item
func (d *DriverDigit) DrawCaptcha(content string) (item Item, err error) {
	// Initialize PRNG.
	itemDigit := NewItemDigit(d.Width, d.Height, d.DotCount, d.MaxSkew)
	// parse digits to string
	digits := string2digits(content)

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
	x := RandomInt(maxX-border*2) + border
	y := RandomInt(maxY-border*2) + border
	// Draw digits.
	for _, n := range digits {
		itemDigit.drawDigit(digitFont[n], x, y)
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

// GenerateCaptcha creates captcha content and answer
func (d *DriverDigit) GenerateCaptcha() (id, q, a string) {
	id = RandomString()
	digits := RandomDigits(d.Length)
	a = digits2String(digits)
	return id, a, a
}
