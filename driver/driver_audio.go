package driver

// DriverAudio captcha config for captcha-engine-audio.
type DriverAudio struct {
	// Length Default number of digits in captcha solution.
	Length int
	// Language possible values for lang are "en", "ja", "ru", "zh".
	Language string
}

// NewDriverAudio creates a driver of audio
func NewDriverAudio(d DriverAudio) *DriverAudio {
	return mergeDriverAudio(d)
}

// DrawCaptcha creates audio captcha item
func (d *DriverAudio) DrawCaptcha(content string) (item Item, err error) {
	digits := string2digits(content)
	audio := newAudio(digits, d.Language)
	return audio, nil
}

// GenerateCaptcha creates id,captcha content and answer
func (d *DriverAudio) GenerateCaptcha() (id, q, a string) {
	id = RandomString()
	digits := RandomDigits(d.Length)
	a = digits2String(digits)
	return id, a, a
}

// mergeDriverAudio merges default driver with given audio driver
func mergeDriverAudio(d DriverAudio) *DriverAudio {
	if d.Length == 0 {
		d.Length = DefaultDriverAudio.Length
	}
	if d.Language == "" {
		d.Language = DefaultDriverAudio.Language
	}
	return &d
}
