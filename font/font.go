package font

import (
	"embed"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const FontPath = "sources"

type Font struct {
	fs embed.FS
}

func NewFont(fs embed.FS) *Font {
	return &Font{fs: fs}
}

// LoadFont load font from file.
func (f *Font) LoadFont(name string) *truetype.Font {
	fontPath := strings.Trim(FontPath, "/")
	fontBytes, err := f.fs.ReadFile(fontPath + "/" + name)
	if err != nil {
		panic(err.Error())
	}
	// font file bytes to trueTypeFont
	trueTypeFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err.Error())
	}
	return trueTypeFont
}

// LoadFonts load fonts from dir.
// make the simple-font(RitaSmith.ttf) the first font of trueTypeFonts.
func (f *Font) LoadFonts(names []string) []*truetype.Font {
	if len(names) == 0 {
		return nil
	}
	fonts := make([]*truetype.Font, 0)
	for _, name := range names {
		fonts = append(fonts, f.LoadFont(name))
	}
	return fonts
}

// LoadAll load all fonts.
func (f *Font) LoadAll() []*truetype.Font {
	return f.LoadFonts([]string{
		"3Dumb.ttf",
		"ApothecaryFont.ttf",
		"Comismsh.ttf",
		"DENNEthree-dee.ttf",
		"DeborahFancyDress.ttf",
		"Flim-Flam.ttf",
		"RitaSmith.ttf",
		"actionj.ttf",
		"chromohv.ttf",
		"wqy-microhei.ttc",
	})
}

// LoadChinese load Chinese font.
func (f *Font) LoadChinese() *truetype.Font {
	return f.LoadFont("wqy-microhei.ttc")
}
