package font

import "embed"

// defaultSource Built-in font storage.
//
//go:embed sources/*.ttf
//go:embed sources/*.ttc
var defaultSource embed.FS

var DefaultSource = NewSource(defaultSource)
