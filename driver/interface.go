package driver

import "io"

// Driver captcha interface for captcha engine to write staff
type Driver interface {
	// DrawCaptcha draws binary item
	DrawCaptcha(content string) (item Item, err error)
	// GenerateIdQuestionAnswer creates rand id, content and answer
	GenerateIdQuestionAnswer() (id, q, a string)
}

// Item is captcha item interface
type Item interface {
	// WriteTo writes to a writer
	WriteTo(w io.Writer) (n int64, err error)
	// EncodeB64string encodes as base64 string
	EncodeB64string() string
}
