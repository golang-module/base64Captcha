package driver

import (
	"image/color"
)

// DriverLetter config for letter driver.
type DriverLetter struct {
	// Width Captcha png width in pixel.
	Width int

	// Height png height in pixel.
	Height int

	// Length random string length.
	Length int

	// NoiseCount text noise count.
	NoiseCount int

	// ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
	ShowLineOptions int

	// Source is a Unicode which is the rand string from.
	Source string

	// BgColor captcha image background color (optional)
	BgColor *color.RGBA
}

// NewDriverLetter creates a letter driver.
func NewDriverLetter(d DriverLetter) *DriverLetter {
	return mergeDriverLetter(d)
}

// DrawCaptcha draws captcha item for letter driver
func (d *DriverLetter) DrawCaptcha(content string) (item Item, err error) {
	driverString := NewDriverString(DriverString{
		Width:           d.Width,
		Height:          d.Height,
		Length:          d.Length,
		NoiseCount:      d.NoiseCount,
		ShowLineOptions: d.ShowLineOptions,
		Source:          d.Source,
		BgColor:         d.BgColor,
	})
	return driverString.DrawCaptcha(content)
}

// GenerateCaptcha generates id, content and answer for letter driver.
func (d *DriverLetter) GenerateCaptcha() (id, content, answer string) {
	driverString := NewDriverString(DriverString{
		Width:           d.Width,
		Height:          d.Height,
		Length:          d.Length,
		NoiseCount:      d.NoiseCount,
		ShowLineOptions: d.ShowLineOptions,
		Source:          d.Source,
		BgColor:         d.BgColor,
	})
	return driverString.GenerateCaptcha()
}

// mergeDriverLetter merges default driver with given letter driver
func mergeDriverLetter(d DriverLetter) *DriverLetter {
	if d.Height == 0 {
		d.Height = DefaultDriverString.Height
	}
	if d.Width == 0 {
		d.Width = DefaultDriverString.Width
	}
	if d.ShowLineOptions == 0 {
		d.ShowLineOptions = DefaultDriverString.ShowLineOptions
	}
	if d.NoiseCount == 0 {
		d.NoiseCount = DefaultDriverString.NoiseCount
	}
	if d.Length == 0 {
		d.Length = DefaultDriverString.Length
	}
	if d.Source == "" {
		d.Source = TxtAlphabet
	}
	return &d
}
