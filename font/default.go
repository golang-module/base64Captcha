package font

import "embed"

// defaultSource Built-in font storage.
//
//go:embed sources/*.ttf
//go:embed sources/*.ttc
var defaultFont embed.FS

var DefaultFont = NewFont(defaultFont)
