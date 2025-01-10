package driver

// DriverAudio captcha config for captcha-engine-audio.
type DriverAudio struct {
	// Length Default number of digits in captcha solution.
	Length int
	// Language possible values for lang are "en", "ja", "ru", "zh".
	Language string
}

// NewDriverAudio creates a driver of audio
func NewDriverAudio(length int, language string) *DriverAudio {
	return &DriverAudio{Length: length, Language: language}
}

// DrawCaptcha creates audio captcha item
func (d *DriverAudio) DrawCaptcha(content string) (item Item, err error) {
	digits := string2digits(content)
	audio := newAudio(digits, d.Language)
	return audio, nil
}

// GenerateIdQuestionAnswer creates id,captcha content and answer
func (d *DriverAudio) GenerateIdQuestionAnswer() (id, q, a string) {
	id = RandomString()
	digits := RandomDigits(d.Length)
	a = digits2String(digits)
	return id, a, a
}
