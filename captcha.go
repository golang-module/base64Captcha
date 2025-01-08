// fork from https://github.com/mojocn/base64Captcha

// Package captcha supports digits, numbers,alphabet, arithmetic, audio and digit-alphabet captcha.
// captcha is used for fast development of restfull APIs, web apps and backend services in Go. give a string identifier to the package, and it returns with a base64-encoding-png-string
package captcha

import "strings"

// Captcha basic information.
type Captcha struct {
	Driver Driver
	Store  Store
}

// NewCaptcha creates a captcha instance from driver and store
func NewCaptcha(driver Driver, store Store) *Captcha {
	return &Captcha{Driver: driver, Store: store}
}

// Generate generates a random id, base64 image string or an error if any
func (c *Captcha) Generate() (id, b64s, answer string, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", "", err
	}
	err = c.Store.Set(id, answer)
	if err != nil {
		return "", "", "", err
	}
	b64s = item.EncodeB64string()
	return
}

// Verify by a given id key and remove the captcha value in store,
// return boolean value.
// if you have multiple captcha instances which share a same store.
// You may want to call `store.Verify` method instead.
func (c *Captcha) Verify(id, answer string, clear bool) (match bool) {
	vv := c.Store.Get(id, clear)
	// fix issue for some redis key-value string value
	vv = strings.TrimSpace(vv)
	return vv == strings.TrimSpace(answer)
}
