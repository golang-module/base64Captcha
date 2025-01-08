package captcha

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_parseDigitsToString(t *testing.T) {
	for i := 1; i < 10; i++ {
		digit := randomDigits(i)
		s := parseDigitsToString(digit)
		if len(s) != i {
			t.Error("failed")
		}
	}
}

func Test_stringToFakeByte(t *testing.T) {
	for i := 1; i < 10; i++ {
		digit := randomDigits(i)
		s := parseDigitsToString(digit)
		if len(s) != i {
			t.Error("failed")
		}
		fb := stringToFakeByte(s)
		if !reflect.DeepEqual(fb, digit) {
			t.Error("failed")
		}
	}
}

func Test_randomDigits(t *testing.T) {
	for i := 1; i < 10; i++ {
		digit := randomDigits(i)
		if len(digit) != i {
			t.Error("failed")
		}

	}
}

func Test_randomBytes(t *testing.T) {
	for i := 1; i < 10; i++ {
		digit := randomBytes(i)
		if len(digit) != i {
			t.Error("failed")
		}
	}
}

func Test_randomBytesMod(t *testing.T) {
	for i := 1; i < 10; i++ {
		digit := randomBytesMod(i, 'c')
		if len(digit) != i {
			t.Error("failed")
		}
	}
}

func Test_itemWriteFile(t *testing.T) {
	// todo:::
}

func Test_pathExists(t *testing.T) {
	td := os.TempDir()
	defer os.RemoveAll(td)
	p := filepath.Join(td, RandomId())
	if pathExists(p) {
		t.Error("failed")
	}
	_ = os.MkdirAll(p, os.ModePerm)

	if !pathExists(p) {
		t.Error("failed")
	}
}
