package driver

import "image/color"

// DefaultDriverAudio is a default audio driver
var DefaultDriverAudio = &DriverAudio{
	Length:   6,
	Language: "en",
}

// DefaultDriverMath is a default math driver
var DefaultDriverMath = &DriverMath{
	Width:           100,
	Height:          32,
	ShowLineOptions: 0,
	NoiseCount:      0,
	Fonts:           []string{"wqy-microhei.ttc"},
	BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}

// DefaultDriverDigit is a default digit driver
var DefaultDriverDigit = &DriverDigit{
	Width:    100,
	Height:   32,
	Length:   6,
	MaxSkew:  0,
	DotCount: 0,
}

// DefaultDriverString is a default string driver
var DefaultDriverString = &DriverString{
	Width:           100,
	Height:          32,
	Length:          6,
	ShowLineOptions: 0,
	NoiseCount:      0,
	Source:          TxtAlphabet,
	Fonts:           []string{"wqy-microhei.ttc"},
	BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}

// DefaultDriverChinese is a default chinese driver
var DefaultDriverChinese = &DriverChinese{
	Width:           100,
	Height:          32,
	Length:          6,
	ShowLineOptions: 0,
	NoiseCount:      0,
	Source:          "设想,你在,处理,消费者,的音,频输,出音,频可,能无,论什,么都,没有,任何,输出,或者,它可,能是,单声道,立体声,或是,环绕立,体声的,,不想要,的值",
	Fonts:           []string{"wqy-microhei.ttc"},
	BgColor:         &color.RGBA{R: 125, G: 125, B: 0, A: 118},
}
