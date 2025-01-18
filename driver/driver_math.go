package driver

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype/truetype"
)

type Fonts struct {
	Name []string
	Data []*truetype.Font
}

// DriverMath config for math driver.
type DriverMath struct {
	// Width Captcha png width in pixel.
	Width int

	// Height png height in pixel.
	Height int

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

// NewDriverMath creates a math driver.
func NewDriverMath(d DriverMath) *DriverMath {
	return mergeDriverMath(d)
}

// DrawCaptcha draws captcha item for math driver.
func (d *DriverMath) DrawCaptcha(question string) (item Item, err error) {
	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = RandomColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	// 波浪线 比较丑
	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	// 背景有文字干扰
	if d.NoiseCount > 0 {
		noise := RandomText(d.NoiseCount, strings.Repeat(TxtNumbers, d.NoiseCount))
		err = itemChar.drawNoise(noise, font.DefaultFont.LoadAll())
		if err != nil {
			return
		}
	}

	// 画 细直线 (n 条)
	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	// 画 多个小波浪线
	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	// draw question
	err = itemChar.drawText(question, d.fontsArray)
	if err != nil {
		return
	}
	return itemChar, nil
}

// GenerateCaptcha generates id, content and answer for math captcha.
func (d *DriverMath) GenerateCaptcha() (id, question, answer string) {
	id = RandomString()
	operators := []string{"+", "-", "x"}
	var result int32
	switch operators[rand.Int31n(3)] {
	case "+":
		a := rand.Int31n(20)
		b := rand.Int31n(20)
		question = fmt.Sprintf("%d+%d=?", a, b)
		result = a + b
	case "x":
		a := rand.Int31n(10)
		b := rand.Int31n(10)
		question = fmt.Sprintf("%dx%d=?", a, b)
		result = a * b
	default:
		a := rand.Int31n(80) + rand.Int31n(20)
		b := rand.Int31n(80)

		question = fmt.Sprintf("%d-%d=?", a, b)
		result = a - b

	}
	answer = fmt.Sprintf("%d", result)
	return
}

// mergeDriverMath merges default driver with given math driver
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
	return &d
}
