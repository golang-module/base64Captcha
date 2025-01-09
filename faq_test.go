package base64Captcha

import (
	"testing"

	"github.com/golang-module/base64Captcha/store"
)

func TestHandlerCaptchaGenerate(t *testing.T) {
	driver := &DriverString{
		Height:          80,
		Width:           240,
		NoiseCount:      10,
		ShowLineOptions: 10,
		Length:          10,
		Source:          "axclajsdlfkjalskjdglasdg",
		BgColor:         nil,
		Fonts:           nil,
	}

	c := NewCaptcha(driver, store.DefaultMemoryStore)

	id, _, _, err := c.Generate()
	if err != nil {
		t.Fatalf("some error: %s", err)
	}

	t.Logf("id: %s", id)
}
