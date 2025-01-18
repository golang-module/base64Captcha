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
	BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}

// DefaultDriverDigit is a default digit driver
var DefaultDriverDigit = &DriverDigit{
	Width:   100,
	Height:  32,
	Length:  6,
	BgColor: &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}

// DefaultDriverLetter is a default letter driver
var DefaultDriverLetter = &DriverLetter{
	Width:           100,
	Height:          32,
	Length:          6,
	ShowLineOptions: 1,
	NoiseCount:      0,
	Source:          TxtAlphabet,
	BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}

// DefaultDriverString is a default string driver
var DefaultDriverString = &DriverString{
	Width:           100,
	Height:          32,
	Length:          6,
	ShowLineOptions: 0,
	NoiseCount:      0,
	Source:          TxtNumbers + TxtAlphabet,
	BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}

// DefaultDriverChinese is a default chinese driver
var DefaultDriverChinese = &DriverChinese{
	Width:           100,
	Height:          32,
	Length:          6,
	ShowLineOptions: 0,
	NoiseCount:      0,
	Source:          "县果栋容他锹射纳堤洲冶架缓飞善挑捏绒既寨剧缝辆语愉谱鸟详坛饶碰扛笔试晶巴呀塘有谣辜确丝活将宪染淋范殖",
	Fonts:           []string{"wqy-microhei.ttc"},
	BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
}
