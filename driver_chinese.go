package base64Captcha

import (
	"image/color"
	"strings"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype/truetype"
)

// DriverChinese is a driver of Unicode Chinese characters.
type DriverChinese struct {
	// Height png height in pixel.
	Height int
	// Width Captcha png width in pixel.
	Width int

	// NoiseCount text noise count.
	NoiseCount int

	// ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
	ShowLineOptions int

	// Length random string length.
	Length int

	// Source is a Unicode which is the rand string from.
	Source string

	// BgColor captcha image background color (optional)
	BgColor *color.RGBA

	// Fonts loads by name see fonts.go's comment
	Fonts      []string
	fontsArray []*truetype.Font
}

// NewDriverChinese creates a driver of Chinese characters
func NewDriverChinese(height int, width int, noiseCount int, showLineOptions int, length int, source string, bgColor *color.RGBA, fonts []string) *DriverChinese {
	defaultSource := font.DefaultSource
	fontsArray := defaultSource.LoadFonts(fonts)
	if len(fontsArray) == 0 {
		fontsArray = defaultSource.LoadAll()
	}
	return &DriverChinese{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, Source: source, BgColor: bgColor, fontsArray: fontsArray}
}

// ConvertFonts loads fonts by names
func (d *DriverChinese) ConvertFonts() *DriverChinese {
	defaultSource := font.DefaultSource
	fontsArray := defaultSource.LoadFonts(d.Fonts)
	if len(fontsArray) == 0 {
		fontsArray = defaultSource.LoadAll()
	}
	d.fontsArray = fontsArray
	return d
}

// GenerateIdQuestionAnswer generates captcha content and its answer
func (d *DriverChinese) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomId()
	ss := strings.Split(d.Source, ",")
	length := len(ss)
	if length == 1 {
		c := RandText(d.Length, ss[0])
		return id, c, c
	}
	if length <= d.Length {
		c := RandText(d.Length, TxtNumbers+TxtAlphabet)
		return id, c, c
	}

	res := make([]string, d.Length)
	for k := range res {
		res[k] = ss[randIntn(length)]
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
		bgc = RandLightColor()
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
		noise := RandText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
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
