package driver

import (
	"image/color"
	"log"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype/truetype"
)

// https://en.wikipedia.org/wiki/Unicode_block
var langMap = map[string][]int{
	// "zh-CN": []int{19968, 40869},
	"latin":  {0x0000, 0x007f},
	"zh":     {0x4e00, 0x9fa5},
	"ko":     {12593, 12686},
	"jp":     {12449, 12531}, // []int{12353, 12435}
	"ru":     {1025, 1169},
	"th":     {0x0e00, 0x0e7f},
	"greek":  {0x0380, 0x03ff},
	"arabic": {0x0600, 0x06ff},
	"hebrew": {0x0590, 0x05ff},
	// "emotion": []int{0x1f601, 0x1f64f},
}

func generateRandomRune(size int, code string) string {
	lang, ok := langMap[code]
	if !ok {
		log.Printf("can not font language of %s \n", code)
		lang = langMap["latin"]
	}
	start := lang[0]
	end := lang[1]
	randRune := make([]rune, size)
	for i := range randRune {
		idx := RandomInt(end-start) + start
		randRune[i] = rune(idx)
	}
	return string(randRune)
}

// DriverLanguage generates language Unicode by language
type DriverLanguage struct {
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

	// BgColor captcha image background color (optional)
	BgColor *color.RGBA

	// Fonts loads by name see fonts.go's comment
	Fonts        []*truetype.Font
	LanguageCode string
}

// NewDriverLanguage creates a driver
func NewDriverLanguage(height int, width int, noiseCount int, showLineOptions int, length int, bgColor *color.RGBA, fonts []*truetype.Font, languageCode string) *DriverLanguage {
	return &DriverLanguage{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, BgColor: bgColor, Fonts: fonts, LanguageCode: languageCode}
}

// GenerateIdQuestionAnswer creates content and answer
func (d *DriverLanguage) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = RandomString()
	content = generateRandomRune(d.Length, d.LanguageCode)
	return id, content, content
}

// DrawCaptcha creates item
func (d *DriverLanguage) DrawCaptcha(content string) (item Item, err error) {
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

	defaultFont := font.DefaultFont

	// draw noise
	if d.NoiseCount > 0 {
		noise := RandomText(d.NoiseCount, TxtNumbers+TxtAlphabet+",.[]<>")
		err = itemChar.DrawNoise(noise, defaultFont.LoadAll())
		if err != nil {
			return
		}
	}

	// draw content
	// use font that match your language
	fontChinese := defaultFont.LoadChinese()
	err = itemChar.DrawText(content, []*truetype.Font{fontChinese})
	if err != nil {
		return
	}

	return itemChar, nil
}
