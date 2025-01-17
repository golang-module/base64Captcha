package driver

import (
	"image/color"
	"strings"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype/truetype"
)

// DriverChinese is a driver of Unicode Chinese characters.
type DriverChinese struct {
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

	// Fonts loads by name see fonts.go's comment
	Fonts      []string
	fontsArray []*truetype.Font
}

// NewDriverChinese creates a driver of Chinese characters
func NewDriverChinese(d DriverChinese) *DriverChinese {
	defaultFont := font.DefaultFont
	fontsArray := defaultFont.LoadFonts(d.Fonts)
	if len(fontsArray) == 0 {
		d.fontsArray = defaultFont.LoadAll()
	}
	return mergeDriverChinese(d)
}

// GenerateCaptcha generates captcha content and its answer
func (d *DriverChinese) GenerateCaptcha() (id, content, answer string) {
	id = RandomString()
	ss := strings.Split(d.Source, ",")
	length := len(ss)
	if length == 1 {
		c := RandomText(d.Length, ss[0])
		return id, c, c
	}
	if length <= d.Length {
		c := RandomText(d.Length, TxtNumbers+TxtAlphabet)
		return id, c, c
	}

	res := make([]string, d.Length)
	for k := range res {
		res[k] = ss[RandomInt(length)]
	}

	content = strings.Join(res, "")
	return id, content, content
}

// DrawCaptcha generates captcha item(image)
func (d *DriverChinese) DrawCaptcha(content string) (item Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandomColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	// draw hollow line
	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	// draw slime line
	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	// draw sine line
	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	// draw noise
	if d.NoiseCount > 0 {
		source := TxtNumbers + TxtAlphabet + ",.[]<>"
		noise := RandomText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
		err = itemChar.drawNoise(noise, d.fontsArray)
		if err != nil {
			return
		}
	}

	// draw content
	err = itemChar.drawText(content, d.fontsArray)
	if err != nil {
		return
	}

	return itemChar, nil
}

// mergeDriverChinese merges default driver with given chinese driver
func mergeDriverChinese(d DriverChinese) *DriverChinese {
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
		d.Source = DefaultDriverString.Source
	}
	if len(d.Fonts) == 0 {
		d.Fonts = DefaultDriverString.Fonts
	}
	if d.BgColor == nil {
		d.BgColor = DefaultDriverString.BgColor
	}
	if len(d.fontsArray) == 0 {
		d.fontsArray = DefaultDriverString.fontsArray
	}
	return &d
}
