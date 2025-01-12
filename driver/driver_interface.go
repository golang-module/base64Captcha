package driver

// Driver captcha interface for captcha engine to write staff
type Driver interface {
	// DrawCaptcha draws binary item
	DrawCaptcha(content string) (item Item, err error)
	// GenerateCaptcha creates rand id, content and answer
	GenerateCaptcha() (id, q, a string)
}
