package driver

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"math"

	"github.com/golang-module/base64Captcha/font"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	imageFont "golang.org/x/image/font"
)

// ItemChar captcha item of unicode characters
type ItemChar struct {
	bgColor color.Color
	width   int
	height  int
	nrgba   *image.NRGBA
}

// NewItemChar creates a captcha item of characters
func NewItemChar(w int, h int, bgColor color.RGBA) *ItemChar {
	d := ItemChar{width: w, height: h}
	m := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(m, m.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)
	d.nrgba = m
	return &d
}

// DrawHollowLine draw strong and bold white line.
func (item *ItemChar) DrawHollowLine() *ItemChar {

	first := item.width / 20
	end := first * 19

	lineColor := RandomColor()

	x1 := float64(RandomInt(first))
	// y1 := float64(RandomInt(y)+y);

	x2 := float64(RandomInt(first) + end)

	multiple := float64(RandomInt(5)+3) / float64(5)
	if int(multiple*10)%3 == 0 {
		multiple = multiple * -1.0
	}

	w := item.height / 20

	for ; x1 < x2; x1++ {

		y := math.Sin(x1*math.Pi*multiple/float64(item.width)) * float64(item.height/3)

		if multiple < 0 {
			y = y + float64(item.height/2)
		}
		item.nrgba.Set(int(x1), int(y), lineColor)

		for i := 0; i <= w; i++ {
			item.nrgba.Set(int(x1), int(y)+i, lineColor)
		}
	}

	return item
}

// DrawSineLine draw a sine line.
func (item *ItemChar) DrawSineLine() *ItemChar {
	var py float64

	// 振幅
	h := RandomInt(item.height / 2)

	// X 轴方向偏移量
	x := RandomRange(float64(-item.height/4), float64(item.height/4))

	// Y 轴方向偏移量
	y := RandomRange(float64(-item.height/4), float64(item.height/4))

	// 周期
	var t float64
	if item.height > item.width/2 {
		t = RandomRange(float64(item.width/2), float64(item.height))
	} else if item.height == item.width/2 {
		t = float64(item.height)
	} else {
		t = RandomRange(float64(item.height), float64(item.width/2))
	}
	w := (2 * math.Pi) / t

	// 曲线横坐标起始位置
	px1 := 0
	px2 := int(RandomRange(float64(item.width)*0.8, float64(item.width)))

	c := RandDeepColor()

	for px := px1; px < px2; px++ {
		if w != 0 {
			py = float64(h)*math.Sin(w*float64(px)+x) + y + (float64(item.width) / float64(5))
			i := item.height / 5
			for i > 0 {
				item.nrgba.Set(px+i, int(py), c)
				i--
			}
		}
	}

	return item
}

// DrawSlimLine draw n slim-random-color lines.
func (item *ItemChar) DrawSlimLine(num int) *ItemChar {

	first := item.width / 10
	end := first * 9

	y := item.height / 3

	for i := 0; i < num; i++ {

		point1 := point{X: RandomInt(first), Y: RandomInt(y)}
		point2 := point{X: RandomInt(first) + end, Y: RandomInt(y)}

		if i%2 == 0 {
			point1.Y = RandomInt(y) + y*2
			point2.Y = RandomInt(y)
		} else {
			point1.Y = RandomInt(y) + y*(i%2)
			point2.Y = RandomInt(y) + y*2
		}

		item.DrawBeeline(point1, point2, RandDeepColor())

	}
	return item
}

// DrawBeeline draw a beeline.
func (item *ItemChar) DrawBeeline(point1 point, point2 point, lineColor color.RGBA) {
	dx := math.Abs(float64(point1.X - point2.X))
	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y >= point2.Y {
		sy = -1
	}
	err := dx - dy
	for {
		item.nrgba.Set(point1.X, point1.Y, lineColor)
		item.nrgba.Set(point1.X+1, point1.Y, lineColor)
		item.nrgba.Set(point1.X-1, point1.Y, lineColor)
		item.nrgba.Set(point1.X+2, point1.Y, lineColor)
		item.nrgba.Set(point1.X-2, point1.Y, lineColor)
		if point1.X == point2.X && point1.Y == point2.Y {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

// DrawNoise draw noise text.
func (item *ItemChar) DrawNoise(noiseText string, fonts []*truetype.Font) error {
	c := freetype.NewContext()
	c.SetDPI(imageStringDpi)

	c.SetClip(item.nrgba.Bounds())
	c.SetDst(item.nrgba)
	c.SetHinting(imageFont.HintingFull)
	rawFontSize := float64(item.height) / (1 + float64(RandomInt(7))/float64(10))

	for _, char := range noiseText {
		rw := RandomInt(item.width)
		rh := RandomInt(item.height)
		fontSize := rawFontSize/2 + float64(RandomInt(5))
		c.SetSrc(image.NewUniform(RandomColor()))
		c.SetFontSize(fontSize)
		c.SetFont(RandomFont(fonts))
		pt := freetype.Pt(rw, rh)
		if _, err := c.DrawString(string(char), pt); err != nil {
			log.Println(err)
		}
	}
	return nil
}

// drawText draw captcha string to image.

func (item *ItemChar) DrawText(text string, fonts []*truetype.Font) error {
	c := freetype.NewContext()
	c.SetDPI(imageStringDpi)
	c.SetClip(item.nrgba.Bounds())
	c.SetDst(item.nrgba)
	c.SetHinting(imageFont.HintingFull)

	if len(text) == 0 {
		return errors.New("text must not be empty, there is nothing to draw")
	}

	fontWidth := item.width / len(text)

	for i, s := range text {
		fontSize := item.height * (RandomInt(7) + 7) / 16
		c.SetSrc(image.NewUniform(RandDeepColor()))
		c.SetFontSize(float64(fontSize))
		c.SetFont(RandomFont(fonts))
		x := fontWidth*i + fontWidth/fontSize
		y := item.height/2 + fontSize/2 - RandomInt(item.height/16*3)
		pt := freetype.Pt(x, y)
		if _, err := c.DrawString(string(s), pt); err != nil {
			return err
		}
	}
	return nil
}

// BinaryEncoding encodes an image to PNG and returns a byte slice.
func (item *ItemChar) BinaryEncoding() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, item.nrgba); err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

// WriteTo writes captcha character in png format into the given io.Writer, and
// returns the number of bytes written and an error if any.
func (item *ItemChar) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(item.BinaryEncoding())
	return int64(n), err
}

// EncodeB64string encodes an image to base64 string
func (item *ItemChar) EncodeB64string() string {
	return fmt.Sprintf("data:%s;base64,%s", MimeTypeImage, base64.StdEncoding.EncodeToString(item.BinaryEncoding()))
}

type point struct {
	X int
	Y int
}

// RandomFont get random font.
func RandomFont(fonts []*truetype.Font) *truetype.Font {
	n := len(fonts)
	if n == 0 {
		// loading default fonts
		fonts = font.DefaultSource.LoadAll()
		n = len(fonts)
	}
	return fonts[RandomInt(n)]
}

// RandDeepColor get random deep color. 随机生成深色系.
func RandDeepColor() color.RGBA {

	randColor := RandColor()

	increase := float64(30 + RandomInt(255))

	red := math.Abs(math.Min(float64(randColor.R)-increase, 255))

	green := math.Abs(math.Min(float64(randColor.G)-increase, 255))
	blue := math.Abs(math.Min(float64(randColor.B)-increase, 255))

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

// RandColor get random color. 生成随机颜色.
func RandColor() color.RGBA {
	red := RandomInt(255)
	green := RandomInt(255)
	var blue int
	if (red + green) > 400 {
		blue = 0
	} else {
		blue = 400 - green - red
	}
	if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}
