package driver

import (
	"io"
)

// Item is captcha item interface
type Item interface {
	// Writer writes to a writer
	Writer(w io.Writer) (n int64, err error)
	// Encoder encodes as base64 string
	Encoder() string
}
