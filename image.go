// Copyright 2011-2014 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
)

const (
	// Standard width and height of a captcha image.
	StdWidth  = 240
	StdHeight = 80
	// Maximum absolute skew factor of a single digit.
	maxSkew = 0.7
	// Number of background circles.
	circleCount = 20
	distortCircleCount = 5
)

type Image struct {
	*image.RGBA
	numWidth  int
	numHeight int
	dotSize   int
	rng       siprng
}

var background color.Color = color.RGBA{144, 144, 144, 0xFF}

// NewImage returns a new captcha image of the given width and height with the
// given digits, where each digit must be in range 0-9.
func NewImage(id string, digits []byte, width, height int) *Image {
	m := new(Image)

	// Initialize PRNG.
	m.rng.Seed(deriveSeed(imageSeedPurpose, id, digits))
	
	prim := color.RGBA{
		uint8(m.rng.Intn(129)),
		uint8(m.rng.Intn(129)),
		uint8(m.rng.Intn(129)),
		0xFF,
	}

	m.RGBA = image.NewRGBA(image.Rect(0, 0, width, height))
	m.calculateSizes(width, height, len(digits))
	// Randomly position captcha inside the image.
	maxx := width - (m.numWidth+m.dotSize)*len(digits) - m.dotSize
	maxy := height - m.numHeight - m.dotSize*2
	var border int
	if width > height {
		border = height / 10
	} else {
		border = width / 10
	}
	x := m.rng.Int(border, maxx-border)
	y := m.rng.Int(border, maxy-border)
	// Draw digits.
	for _, n := range digits {
		m.drawDigit(getChar(n), x, y, prim)
		x += m.numWidth + m.dotSize
	}
	// Draw strike-through line.
	m.strikeThrough(prim)
	
	m.fillWithDistortCircles(distortCircleCount, width/4)
	// Apply wave distortion.
	m.distort(m.rng.Float(5, 10), m.rng.Float(100, 200))
	// Fill image with random circles.
	m.fillWithCircles(circleCount, m.dotSize)
	return m
}

/*
func (m *Image) getRandomPalette() color.Palette {
	p := make([]color.Color, circleCount+1)
	// Transparent color.
	p[0] = color.RGBA{0xFF, 0xFF, 0xFF, 0x00}
	// Primary color.
	prim := color.RGBA{
		uint8(m.rng.Intn(129)),
		uint8(m.rng.Intn(129)),
		uint8(m.rng.Intn(129)),
		0xFF,
	}
	p[1] = prim
	// Circle colors.
	for i := 2; i <= circleCount; i++ {
		p[i] = m.randomBrightness(prim, 255)
	}
	return p
}*/

// encodedPNG encodes an image to PNG and returns
// the result as a byte slice.
func (m *Image) encodedPNG() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		panic(err.Error())
	}
	return buf.Bytes()
}

// WriteTo writes captcha image in PNG format into the given writer.
func (m *Image) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m.encodedPNG())
	return int64(n), err
}

func (m *Image) calculateSizes(width, height, ncount int) {
	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 8
	} else {
		border = width / 8
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)
	h := float64(height - border*2)
	// fw takes into account 1-dot spacing between digits.
	fw := float64(fontWidth + 1)
	fh := float64(fontHeight)
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
	m.dotSize = int((nh / fh))
	if m.dotSize < 1 {
		m.dotSize = 1
	}
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	m.numWidth = int(nw) - m.dotSize
	m.numHeight = int(nh)
}

func (m *Image) drawHorizLine(fromX, toX, y int, color color.Color) {
	for x := fromX; x <= toX; x++ {
		m.Set(x, y, color)
	}
}

func (m *Image) drawCircle(x, y, radius int, color color.Color) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xo := 0
	yo := radius

	m.Set(x, y+radius, color)
	m.Set(x, y-radius, color)
	m.drawHorizLine(x-radius, x+radius, y, color)

	for xo < yo {
		if f >= 0 {
			yo--
			dfy += 2
			f += dfy
		}
		xo++
		dfx += 2
		f += dfx
		m.drawHorizLine(x-xo, x+xo, y+yo, color)
		m.drawHorizLine(x-xo, x+xo, y-yo, color)
		m.drawHorizLine(x-yo, x+yo, y+xo, color)
		m.drawHorizLine(x-yo, x+yo, y-xo, color)
	}
}


func (m *Image) colorDistortCircle(xc, yc, r int) {
	d := 2*r
	r2 := r*r
	
	dg := int(m.rng.Intn(240)) - 120
	if dg < 0 {
		dg -= 20
	} else {
		dg += 20
	}
	
	dr := int(m.rng.Intn(240)) - 120
	if dr < 0 {
		dr -= 20
	} else {
		dr += 20
	}
	
	db := int(m.rng.Intn(240)) - 120
	if db < 0 {
		db -= 20
	} else {
		db += 20
	}
	
	for y := 0; y <= d; y++ {
		for x := 0; x <= d; x++ {
			x2 := (x-r)*(x-r)
			y2 := (y-r)*(y-r)
			
			if x2+y2 <= r2 {
				xp := xc - r + x
				yp := yc - r + y
				r,g,b,a := m.At(xp, yp).RGBA()
				if a != 0x00 { //background
					m.Set(xp, yp, color.RGBA{
						uint8(int(r)+dr),
						uint8(int(g)+dg),
						uint8(int(b)+db),
						0xFF,
					})
				}
			}
		}
	}
}


func (m *Image) fillWithCircles(n, maxradius int) {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		color := color.RGBA{
			uint8(m.rng.Intn(129)),
			uint8(m.rng.Intn(129)),
			uint8(m.rng.Intn(129)),
			0xFF,
		}
		r := m.rng.Int(1, maxradius)
		m.drawCircle(m.rng.Int(r, maxx-r), m.rng.Int(r, maxy-r), r, color)
	}
}

func (m *Image) fillWithDistortCircles(n, maxradius int) {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	for i := 0; i < n; i++ {
		r := m.rng.Int(maxradius/10, maxradius-1)
		m.colorDistortCircle(m.rng.Int(0, maxx), m.rng.Int(0, maxy), r)
	}
}

func (m *Image) strikeThrough(color color.Color) {
	maxx := m.Bounds().Max.X
	maxy := m.Bounds().Max.Y
	y := m.rng.Int(maxy/3, maxy-maxy/3)
	amplitude := m.rng.Float(5, 20)
	period := m.rng.Float(80, 180)
	dx := 2.0 * math.Pi / period
	for x := 0; x < maxx; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		for yn := 0; yn < m.dotSize; yn++ {
			r := m.rng.Int(0, m.dotSize)
			m.drawCircle(x+int(xo), y+int(yo)+(yn*m.dotSize), r/2, color)
		}
	}
}

func (m *Image) drawDigit(digit []byte, x, y int, color color.Color) {
	skf := m.rng.Float(-maxSkew, maxSkew)
	xs := float64(x)
	r := m.dotSize / 2
	y += m.rng.Int(-r, r)
	for yo := 0; yo < fontHeight; yo++ {
		for xo := 0; xo < fontWidth; xo++ {
			if digit[yo*fontWidth+xo] != blackChar {
				continue
			}
			m.drawCircle(x+xo*m.dotSize, y+yo*m.dotSize, r, color)
		}
		xs += skf
		x = int(xs)
	}
}

func (m *Image) distort(amplude float64, period float64) {
	w := m.Bounds().Max.X
	h := m.Bounds().Max.Y
	r32,g32,b32,_ := background.RGBA()
	r:=uint8(r32)
	g:=uint8(g32)
	b:=uint8(b32)
	dbg := 32.0/float32(w);

	oldm := m
	newm := image.NewRGBA(image.Rect(0, 0, w, h))

	dx := 2.0 * math.Pi / period
	for x := 0; x < w; x++ {
		dbgx:=uint8(float32(x)*dbg)
		bg:=color.RGBA{r-dbgx, g-dbgx, b-dbgx, 0xFF}
		for y := 0; y < h; y++ {
			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)
			c := oldm.At(x+int(xo), y+int(yo))
			newm.Set(x, y, c)
			_,_,_, a := c.RGBA()
			if a == 0x00 {
				newm.Set(x, y, bg)
			}
		}
	}
	m.RGBA = newm
}

func (m *Image) randomBrightness(c color.RGBA, max uint8) color.RGBA {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return c
	}
	n := m.rng.Intn(int(max-maxc)) - int(minc)
	return color.RGBA{
		uint8(int(c.R) + n),
		uint8(int(c.G) + n),
		uint8(int(c.B) + n),
		uint8(c.A),
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
