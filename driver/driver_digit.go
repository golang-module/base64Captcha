package driver

import (
	"image/color"
)

// DriverDigit config for digit driver.
type DriverDigit struct {
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

// NewDriverDigit creates a digit driver.
func NewDriverDigit(d DriverDigit) *DriverDigit {
	return mergeDriverDigit(d)
}

// DrawCaptcha draws captcha item for digit driver.
func (d *DriverDigit) DrawCaptcha(content string) (item Item, err error) {
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

// GenerateCaptcha generates id, content and answer for digit driver.
func (d *DriverDigit) GenerateCaptcha() (id, content, answer string) {
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

// mergeDriverDigit merges default driver with given digit driver.
func mergeDriverDigit(d DriverDigit) *DriverDigit {
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
		d.Source = TxtNumbers
	}
	return &d
}
