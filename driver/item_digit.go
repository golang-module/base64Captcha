package driver

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
)

const (
	digitFontWidth     = 11
	digitFontHeight    = 18
	digitFontBlackChar = 1
)

// ItemDigit digits captcha Struct
type ItemDigit struct {
	width  int
	height int
	*image.Paletted
	dotSize  int
	dotCount int
	maxSkew  float64
}

// NewItemDigit create an instance of item-digit
func NewItemDigit(width int, height int, dotCount int, maxSkew float64) *ItemDigit {
	itemDigit := &ItemDigit{width: width, height: height, dotCount: dotCount, maxSkew: maxSkew}
	// init image.Paletted
	itemDigit.Paletted = image.NewPaletted(image.Rect(0, 0, width, height), createRandPaletteColors(dotCount))
	return itemDigit
}

func createRandPaletteColors(dotCount int) color.Palette {
	p := make([]color.Color, dotCount+1)
	// Transparent color.
	p[0] = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF}
	// Primary color.
	prim := color.RGBA{
		R: uint8(RandomInt(129)),
		G: uint8(RandomInt(129)),
		B: uint8(RandomInt(129)),
		A: 0xFF,
	}

	if dotCount == 0 {
		p[0] = prim
		return p
	}

	p[1] = prim
	// Circle colors.
	for i := 2; i <= dotCount; i++ {
		p[i] = randomBrightness(prim, 255)
	}
	return p
}

func (m *ItemDigit) calculateSizes(width, height, ncount int) {
	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 4
	} else {
		border = width / 4
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)
	h := float64(height - border*2)
	// fw takes into account 1-dot spacing between digits.
	fw := float64(digitFontWidth + 1)
	fh := float64(digitFontHeight)
	nc := float64(ncount)
	// Calculate the width of a single digit taking into account only the
	// width of the image.
	nw := w / nc
	// Calculate the height of a digit from this width.
	nh := nw * fh / fw
	// Digit too high?
	if nh > h {
		// Fit digits based on height.
		nh = h
		nw = fw / fh * nh
	}
	// Calculate dot size.
	m.dotSize = int(nh / fh)
	if m.dotSize < 1 {
		m.dotSize = 1
	}
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	m.width = int(nw) - m.dotSize
	m.height = int(nh)
}

func (m *ItemDigit) drawHorizLine(fromX, toX, y int, colorIdx uint8) {
	for x := fromX; x <= toX; x++ {
		m.SetColorIndex(x, y, colorIdx)
	}
}

func (m *ItemDigit) drawCircle(x, y, radius int, colorIdx uint8) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xo := 0
	yo := radius

	m.SetColorIndex(x, y+radius, colorIdx)
	m.SetColorIndex(x, y-radius, colorIdx)
	m.drawHorizLine(x-radius, x+radius, y, colorIdx)

	for xo < yo {
		if f >= 0 {
			yo--
			dfy += 2
			f += dfy
		}
		xo++
		dfx += 2
		f += dfx
		m.drawHorizLine(x-xo, x+xo, y+yo, colorIdx)
		m.drawHorizLine(x-xo, x+xo, y-yo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y+xo, colorIdx)
		m.drawHorizLine(x-yo, x+yo, y-xo, colorIdx)
	}
}

func (m *ItemDigit) fillWithCircles(n, maxradius int) {
	maxX := m.Bounds().Max.X
	maxY := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		// colorIdx := uint8(m.rng.Int(1, m.dotCount-1))
		colorIdx := uint8(RandomRange(1, m.dotCount-1))
		// r := m.rng.Int(1, maxradius)
		r := RandomRange(1, maxradius)
		// m.drawCircle(m.rng.Int(r, maxx-r), m.rng.Int(r, maxy-r), r, colorIdx)
		m.drawCircle(RandomRange(r, maxX-r), RandomRange(r, maxY-r), r, colorIdx)
	}
}

func (m *ItemDigit) strikeThrough() {
	maxX := m.Bounds().Max.X
	maxY := m.Bounds().Max.Y
	y := RandomRange(maxY/3, maxY-maxY/3)
	amplitude := RandomRange(5.0, 20.0)
	period := RandomRange(80.0, 180.0)
	dx := 2.0 * math.Pi / period
	for x := 0; x < maxX; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		for yn := 0; yn < m.dotSize; yn++ {
			// r := m.rng.Int(0, m.dotSize)
			r := RandomInt(m.dotSize)
			m.drawCircle(x+int(xo), y+int(yo)+(yn*m.dotSize), r/2, 1)
		}
	}
}

// draw digit
func (m *ItemDigit) drawDigit(digit []byte, x, y int) {
	skf := RandomRange(-m.maxSkew, m.maxSkew)
	xs := float64(x)
	r := m.dotSize / 2
	y += RandomRange(-r, r)
	for yo := 0; yo < digitFontHeight; yo++ {
		for xo := 0; xo < digitFontWidth; xo++ {
			if digit[yo*digitFontWidth+xo] != digitFontBlackChar {
				continue
			}
			m.drawCircle(x+xo*m.dotSize, y+yo*m.dotSize, r, 1)
		}
		xs += skf
		x = int(xs)
	}
}

func (m *ItemDigit) distort(amplude float64, period float64) {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y

	oldItem := m.Paletted
	newItem := image.NewPaletted(image.Rect(0, 0, w, h), oldItem.Palette)

	dx := 2.0 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			newItem.SetColorIndex(x, y, oldItem.ColorIndexAt(x+int(xo), y+int(yo)))
		}
	}
	m.Paletted = newItem
}

func randomBrightness(c color.RGBA, max uint8) color.RGBA {
	minC := min3(c.R, c.G, c.B)
	maxC := max3(c.R, c.G, c.B)
	if maxC > max {
		return c
	}
	n := RandomInt(int(max-maxC)) - int(minC)
	return color.RGBA{
		R: uint8(int(c.R) + n),
		G: uint8(int(c.G) + n),
		B: uint8(int(c.B) + n),
		A: c.A,
	}
}

func min3(x, y, z uint8) (m uint8) {
	m = x
	if y < m {
		m = y
	}
	if z < m {
		m = z
	}
	return
}

func max3(x, y, z uint8) (m uint8) {
	m = x
	if y > m {
		m = y
	}
	if z > m {
		m = z
	}
	return
}

// encodeBinary encodes an image to PNG and returns a byte slice.
func (m *ItemDigit) encodeBinary() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m.Paletted); err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

// Writer writes captcha character in png format into the given io.Writer, and
// returns the number of bytes written and an error if any.
func (m *ItemDigit) Writer(w io.Writer) (int64, error) {
	n, err := w.Write(m.encodeBinary())
	return int64(n), err
}

// Encoder encodes an image to base64 string
func (m *ItemDigit) Encoder() string {
	return fmt.Sprintf("data:%s;base64,%s", MimeTypeImage, base64.StdEncoding.EncodeToString(m.encodeBinary()))
}

// converts string to digits
func string2digits(content string) []byte {
	digits := make([]byte, len(content))
	for idx, cc := range content {
		digits[idx] = byte(cc - '0')
	}
	return digits
}

// converts digits to string
func digits2String(bytes []byte) string {
	stringB := make([]byte, len(bytes))
	for idx, by := range bytes {
		stringB[idx] = by + '0'
	}
	return string(stringB)
}
