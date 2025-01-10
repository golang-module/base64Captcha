package base64Captcha

import (
	mathRand "math/rand"
	"reflect"
	"testing"

	"github.com/golang-module/base64Captcha/driver"
	"github.com/golang-module/base64Captcha/store"
)

func TestCaptcha_GenerateB64s(t *testing.T) {
	type fields struct {
		Driver driver.Driver
		Store  store.Store
	}

	dDigit := driver.DriverDigit{Height: 80, Width: 240, Length: 5, MaxSkew: 0.7, DotCount: 5}
	audioDriver := driver.NewDriverAudio(driver.DriverAudio{
		Length:   randIntn(5),
		Language: "en",
	})
	tests := []struct {
		name     string
		fields   fields
		wantId   string
		wantB64s string
		wantErr  bool
	}{
		{"mem-digit", fields{&dDigit, store.DefaultStoreMemory}, "xxxx", "", false},
		{"mem-audio", fields{audioDriver, store.DefaultStoreMemory}, "xxxx", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCaptcha(tt.fields.Driver, tt.fields.Store)
			gotId, b64s, _, err := c.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Captcha.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(b64s)

			a := c.Store.Get(gotId, false)
			if !c.Verify(gotId, a, true) {
				t.Error("false")
			}
		})
	}
}

func TestCaptcha_Verify(t *testing.T) {
	type fields struct {
		Driver driver.Driver
		Store  store.Store
	}
	type args struct {
		id     string
		answer string
		clear  bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantMatch bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Captcha{
				Driver: tt.fields.Driver,
				Store:  tt.fields.Store,
			}
			if gotMatch := c.Verify(tt.args.id, tt.args.answer, tt.args.clear); gotMatch != tt.wantMatch {
				t.Errorf("Captcha.Verify() = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestNewCaptcha(t *testing.T) {
	type args struct {
		driver driver.Driver
		store  store.Store
	}
	tests := []struct {
		name string
		args args
		want *Captcha
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCaptcha(tt.args.driver, tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCaptcha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCaptcha_Generate(t *testing.T) {
	tests := []struct {
		name     string
		c        *Captcha
		wantId   string
		wantB64s string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotB64s, _, err := tt.c.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Captcha.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("Captcha.Generate() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotB64s != tt.wantB64s {
				t.Errorf("Captcha.Generate() gotB64s = %v, want %v", gotB64s, tt.wantB64s)
			}
		})
	}
}

func randIntn(n int) int {
	if n > 0 {
		return mathRand.Intn(n)
	}
	return 0
}
