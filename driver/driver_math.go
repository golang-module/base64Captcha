package driver

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype/truetype"
)

// DriverMath captcha config for captcha math
type DriverMath struct {
	// Height png height in pixel.
	Height int

	// Width Captcha png width in pixel.
	Width int

	// NoiseCount text noise count.
	NoiseCount int

	// ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
	ShowLineOptions int

	// BgColor captcha image background color (optional)
	BgColor *color.RGBA

	// Fonts loads by name see fonts.go's comment
	Fonts      []string
	fontsArray []*truetype.Font
}

// NewDriverMath creates a driver of math
func NewDriverMath(d DriverMath) *DriverMath {
	defaultFont := font.DefaultFont
	fontsArray := defaultFont.LoadFonts(d.Fonts)
	if len(fontsArray) == 0 {
		d.fontsArray = defaultFont.LoadAll()
	}
	return mergeDriverMath(d)
}

// mergeDriverMath merge default config
func mergeDriverMath(d DriverMath) *DriverMath {
	if d.Height == 0 {
		d.Height = DefaultDriverMath.Height
	}
	if d.Width == 0 {
		d.Width = DefaultDriverMath.Width
	}
	if d.ShowLineOptions == 0 {
		d.ShowLineOptions = DefaultDriverMath.ShowLineOptions
	}
	if d.NoiseCount == 0 {
		d.NoiseCount = DefaultDriverMath.NoiseCount
	}
	if len(d.Fonts) == 0 {
		d.Fonts = DefaultDriverMath.Fonts
	}
	if d.BgColor == nil {
		d.BgColor = DefaultDriverMath.BgColor
	}
	if len(d.fontsArray) == 0 {
		d.fontsArray = DefaultDriverMath.fontsArray
	}
	return &d
}

// ConvertFonts loads fonts from names
func (d *DriverMath) ConvertFonts() *DriverMath {
	defaultFont := font.DefaultFont
	fontsArray := defaultFont.LoadFonts(d.Fonts)
	if len(fontsArray) == 0 {
		fontsArray = defaultFont.LoadAll()
	}
	return d
}

// GenerateIdQuestionAnswer creates id,captcha content and answer
func (d *DriverMath) GenerateIdQuestionAnswer() (id, question, answer string) {
	id = RandomString()
	operators := []string{"+", "-", "x"}
	var mathResult int32
	switch operators[rand.Int31n(3)] {
	case "+":
		a := rand.Int31n(20)
		b := rand.Int31n(20)
		question = fmt.Sprintf("%d+%d=?", a, b)
		mathResult = a + b
	case "x":
		a := rand.Int31n(10)
		b := rand.Int31n(10)
		question = fmt.Sprintf("%dx%d=?", a, b)
		mathResult = a * b
	default:
		a := rand.Int31n(80) + rand.Int31n(20)
		b := rand.Int31n(80)

		question = fmt.Sprintf("%d-%d=?", a, b)
		mathResult = a - b

	}
	answer = fmt.Sprintf("%d", mathResult)
	return
}

// DrawCaptcha creates math captcha item
func (d *DriverMath) DrawCaptcha(question string) (item Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	// 波浪线 比较丑
	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.DrawHollowLine()
	}

	// 背景有文字干扰
	if d.NoiseCount > 0 {
		noise := RandomText(d.NoiseCount, strings.Repeat(TxtNumbers, d.NoiseCount))
		err = itemChar.DrawNoise(noise, font.DefaultFont.LoadAll())
		if err != nil {
			return
		}
	}

	// 画 细直线 (n 条)
	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.DrawSlimLine(3)
	}

	// 画 多个小波浪线
	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.DrawSineLine()
	}

	// draw question
	err = itemChar.DrawText(question, d.fontsArray)
	if err != nil {
		return
	}
	return itemChar, nil
}
