package font

import (
	"embed"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const FontPath = "sources"

type Source struct {
	fs embed.FS
}

func NewSource(fs embed.FS) *Source {
	return &Source{fs: fs}
}

// LoadFont import font from file.
func (s *Source) LoadFont(name string) *truetype.Font {
	fontPath := strings.TrimLeft(FontPath, "/")
	fontBytes, err := s.fs.ReadFile(fontPath + "/" + name)
	if err != nil {
		panic(err)
	}
	// font file bytes to trueTypeFont
	trueTypeFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	return trueTypeFont
}

// LoadFonts import fonts from dir.
// make the simple-font(RitaSmith.ttf) the first font of trueTypeFonts.
func (s *Source) LoadFonts(names []string) []*truetype.Font {
	fonts := make([]*truetype.Font, 0)
	for _, name := range names {
		f := s.LoadFont(name)
		fonts = append(fonts, f)
	}
	return fonts
}

// LoadAll import all fonts.
func (s *Source) LoadAll() []*truetype.Font {
	return s.LoadFonts([]string{
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

// LoadChinese import Chinese font.
func (s *Source) LoadChinese() *truetype.Font {
	return s.LoadFont("wqy-microhei.ttc")
}
