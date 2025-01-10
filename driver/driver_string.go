package driver

import (
	"image/color"
	"strings"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype/truetype"
)

// DriverString captcha config for captcha-engine-characters.
type DriverString struct {
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

// NewDriverString creates driver
func NewDriverString(height int, width int, noiseCount int, showLineOptions int, length int, source string, bgColor *color.RGBA, fonts []string) *DriverString {
	defaultFont := font.DefaultFont
	fontsArray := defaultFont.LoadFonts(fonts)
	if len(fontsArray) == 0 {
		fontsArray = defaultFont.LoadAll()
	}
	return &DriverString{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, Source: source, BgColor: bgColor, fontsArray: fontsArray, Fonts: fonts}
}

// ConvertFonts loads fonts by names
func (d *DriverString) ConvertFonts() *DriverString {
	defaultFont := font.DefaultFont
	fontsArray := defaultFont.LoadFonts(d.Fonts)
	if len(fontsArray) == 0 {
		fontsArray = defaultFont.LoadAll()
	}
	d.fontsArray = fontsArray
	return d
}

// GenerateIdQuestionAnswer creates id,content and answer
func (d *DriverString) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomString()
	content = RandomText(d.Length, d.Source)
	return id, content, content
}

// DrawCaptcha draws captcha item
func (d *DriverString) DrawCaptcha(content string) (item Item, err error) {

	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	// draw hollow line
	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.DrawHollowLine()
	}

	// draw slime line
	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.DrawSlimLine(3)
	}

	// draw sine line
	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.DrawSineLine()
	}

	// draw noise
	if d.NoiseCount > 0 {
		source := TxtNumbers + TxtAlphabet + ",.[]<>"
		noise := RandomText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
		err = itemChar.DrawNoise(noise, d.fontsArray)
		if err != nil {
			return
		}
	}

	// draw content
	err = itemChar.DrawText(content, d.fontsArray)
	if err != nil {
		return
	}

	return itemChar, nil
}
