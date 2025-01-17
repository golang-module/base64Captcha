// fork from https://github.com/mojocn/base64Captcha

// Package base64Captcha supports digits, numbers,alphabet, arithmetic, audio and digit-alphabet captcha.
// base64Captcha is used for fast development of restfull APIs, web apps and backend services in Go. give a string identifier to the package, and it returns with a base64-encoding-png-string
package base64Captcha

import (
	"strings"

	"github.com/golang-module/base64Captcha/driver"
	"github.com/golang-module/base64Captcha/store"
)

// Version current version
const Version = "1.3.9"

// Captcha basic information.
type Captcha struct {
	Driver driver.Driver
	Store  store.Store
}

// NewCaptcha creates a captcha instance from driver and store
func NewCaptcha(d driver.Driver, s ...store.Store) *Captcha {
	if len(s) == 0 {
		s = make([]store.Store, 1)
		s[0] = store.DefaultStoreMemory
	}
	return &Captcha{Driver: d, Store: s[0]}
}

// Generate generates a random id, base64 image string or an error if any
func (c *Captcha) Generate() (id, src, answer string, err error) {
	id, content, answer := c.Driver.GenerateCaptcha()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return
	}
	err = c.Store.Set(id, answer)
	if err != nil {
		return
	}
	src = item.Encoder()
	return
}

// Verify by a given id key and remove the captcha value in store,
// return boolean value.
// if you have multiple captcha instances which share a same store.
// You may want to call `store.Verify` method instead.
func (c *Captcha) Verify(id, answer string, clear bool) (match bool) {
	value := c.Store.Get(id, clear)
	// fix issue for some redis key-value string value
	return strings.EqualFold(strings.TrimSpace(value), strings.TrimSpace(answer))
}
