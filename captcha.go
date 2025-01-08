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
func NewCaptcha(driver Driver, store ...Store) *Captcha {
	if len(store) == 0 {
		store = make([]Store, 1)
		store[0] = DefaultMemStore
	}
	return &Captcha{Driver: driver, Store: store[0]}
}

// Generate generates a random id, base64 image string or an error if any
func (c *Captcha) Generate() (id, src, answer string, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return
	}
	err = c.Store.Set(id, answer)
	if err != nil {
		return
	}
	src = item.EncodeB64string()
	return
}

// Verify by a given id key and remove the captcha value in store,
// return boolean value.
// if you have multiple captcha instances which share a same store.
// You may want to call `store.Verify` method instead.
func (c *Captcha) Verify(id, answer string, clear bool) (match bool) {
	value := c.Store.Get(id, clear)
	// fix issue for some redis key-value string value
	return strings.TrimSpace(value) == strings.TrimSpace(answer)
}
