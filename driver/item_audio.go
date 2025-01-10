package driver

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
)

var endingBeepSound []byte

func init() {
	endingBeepSound = changeSpeed(beepSound, 1.4)
}

// ItemAudio captcha-audio-engine return type.
type ItemAudio struct {
	answer      string
	body        *bytes.Buffer
	digitSounds [][]byte
}

// newAudio returns a new audio captcha with the given digits, where each digit
// must be in range 0-9. Digits are pronounced in the given language. If there
// are no sounds for the given language, English is used.
// Possible values for lang are "en", "ja", "ru", "zh".
func newAudio(digits []byte, lang string) *ItemAudio {
	a := new(ItemAudio)

	if sounds, ok := digitSounds[lang]; ok {
		a.digitSounds = sounds
	} else {
		a.digitSounds = digitSounds["en"]
	}
	numsnd := make([][]byte, len(digits))
	for i, n := range digits {
		snd := a.randomDigitSound(n)
		a.setSoundLevel(snd, 1.5)
		numsnd[i] = snd
	}
	// Random intervals between digits (including beginning).
	intervals := make([]int, len(digits)+1)
	intdur := 0
	for i := range intervals {
		dur := RandomRange(sampleRate, sampleRate*2) // 1 to 2 seconds
		intdur += dur
		intervals[i] = dur
	}
	// Generate background sound.
	bg := a.makeBackgroundSound(a.longestDigitSndLen()*len(digits) + intdur)
	// Create buffer and write audio to it.
	sil := a.makeSilence(sampleRate / 5)
	bufCap := 3*len(beepSound) + 2*len(sil) + len(bg) + len(endingBeepSound)
	a.body = bytes.NewBuffer(make([]byte, 0, bufCap))
	// Write prelude, three beeps.
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	a.body.Write(sil)
	a.body.Write(beepSound)
	// Write digits.
	pos := intervals[0]
	for i, v := range numsnd {
		a.mixSound(bg[pos:], v)
		pos += len(v) + intervals[i+1]
	}
	a.body.Write(bg)
	// Write ending (one beep).
	a.body.Write(endingBeepSound)
	return a
}

func (a *ItemAudio) longestDigitSndLen() int {
	n := 0
	for _, v := range a.digitSounds {
		if n < len(v) {
			n = len(v)
		}
	}
	return n
}

// WriteTo writes captcha audio in WAVE format into the given io.Writer, and
// returns the number of bytes written and an error if any.
func (a *ItemAudio) WriteTo(w io.Writer) (n int64, err error) {
	// Calculate padded length of PCM chunk data.
	bodyLen := uint32(a.body.Len())
	paddedBodyLen := bodyLen
	if bodyLen%2 != 0 {
		paddedBodyLen++
	}
	totalLen := uint32(len(waveHeader)) - 4 + paddedBodyLen
	// Header.
	header := make([]byte, len(waveHeader)+4) // includes 4 bytes for chunk size
	copy(header, waveHeader)
	// Put the length of whole RIFF chunk.
	binary.LittleEndian.PutUint32(header[4:], totalLen)
	// Put the length of WAVE chunk.
	binary.LittleEndian.PutUint32(header[len(waveHeader):], bodyLen)
	// Write header.
	nn, err := w.Write(header)
	n = int64(nn)
	if err != nil {
		return
	}
	// Write data.
	n, err = a.body.WriteTo(w)
	n += int64(nn)
	if err != nil {
		return
	}
	// Pad byte if chunk length is odd.
	// (As header has even length, we can check if n is odd, not chunk).
	if bodyLen != paddedBodyLen {
		w.Write([]byte{0})
		n++
	}
	return
}

// EncodeB64string encodes a sound to base64 string
func (a *ItemAudio) EncodeB64string() string {
	var buf bytes.Buffer
	if _, err := a.WriteTo(&buf); err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("data:%s;base64,%s", MimeTypeAudio, base64.StdEncoding.EncodeToString(buf.Bytes()))
}

// setSoundLevel sets the level of the sound.
func (a *ItemAudio) setSoundLevel(b []byte, level float64) {
	for i, v := range b {
		f := float64(v)
		switch {
		case f > 128:
			if f = (f-128)*level + 128; f < 128 {
				f = 128
			}
		case f < 128:
			if f = 128 - (128-f)*level; f > 128 {
				f = 128
			}
		default:
			continue
		}
		b[i] = byte(f)
	}
}

// reversedSound reverses the given sound.
func (a *ItemAudio) reversedSound(b []byte) []byte {
	n := len(b)
	bs := make([]byte, n)
	for i, v := range b {
		bs[n-1-i] = v
	}
	return bs
}

// mixSound mixes the given sounds.
func (a *ItemAudio) mixSound(dst, src []byte) {
	for i, v := range src {
		s, d := int(v), int(dst[i])
		if s < 128 && d < 128 {
			dst[i] = byte(s * d / 128)
		} else {
			dst[i] = byte(2*(s+d) - s*d/128 - 256)
		}
	}
}

// randomDigitSound returns a digit sound with a random speed.
func (a *ItemAudio) randomDigitSound(n byte) []byte {
	s := a.randomSpeed(a.digitSounds[n])
	a.setSoundLevel(s, RandomRange(0.85, 1.2))
	return s
}

// randomSpeed returns a speed-changed sound.
func (a *ItemAudio) randomSpeed(b []byte) []byte {
	return changeSpeed(b, RandomRange(0.95, 1.1))
}

// changeSpeed changes the speed of the given sound.
func changeSpeed(b []byte, speed float64) []byte {
	r := make([]byte, int(math.Floor(float64(len(b))*speed)))
	var f float64
	for _, v := range b {
		for i := int(f); i < int(f+speed); i++ {
			r[i] = v
		}
		f += speed
	}
	return r
}

// makeSilence returns n silent bytes.
func (a *ItemAudio) makeSilence(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = 128
	}
	return b
}

// makeWhiteNoise returns n white noise bytes.
func (a *ItemAudio) makeWhiteNoise(length int, level uint8) []byte {
	noise := a.makeNoise(length)
	adj := 128 - level/2
	for i, v := range noise {
		v %= level
		v += adj
		noise[i] = v
	}
	return noise
}

// makeNoise returns n random bytes.
func (a *ItemAudio) makeNoise(n int) []byte {
	// Since we don't have a buffer for generated bytes in siprng state,
	// we just generate enough 8-byte blocks and then cut the result to the
	// required length. Doing it this way, we lose generated bytes, and we
	// don't get the strictly sequential deterministic output from PRNG:
	// calling Uint64() and then Bytes(3) produces different output than
	// when calling them in the reverse order, but for our applications
	// this is OK.
	numBlocks := (n + 8 - 1) / 8
	b := make([]byte, numBlocks*8)
	for i := 0; i < len(b); i += 8 {
		binary.LittleEndian.PutUint64(b[i:], rand.Uint64())
	}
	return b[:n]
}

// makeBackgroundSound returns a background sound.
func (a *ItemAudio) makeBackgroundSound(length int) []byte {
	noise := a.makeWhiteNoise(length, 4)
	for i := 0; i < length/(sampleRate/10); i++ {
		sound := a.reversedSound(a.digitSounds[RandomInt(10)])
		// snd = changeSpeed(snd, a.rng.Float(0.8, 1.2))
		place := RandomInt(len(noise) - len(sound))
		a.setSoundLevel(sound, RandomRange(0.04, 0.08))
		a.mixSound(noise[place:], sound)
	}
	return noise
}
