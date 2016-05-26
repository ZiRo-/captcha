// Copyright 2016 ZiRo. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package libgocaptcha

import (
	"encoding/gob"
	"os"
)

const (
	fontWidth  = 15
	fontHeight = 25
	blackChar  = 1
)

type Font struct {
	data map[rune][]byte
}

var fonts = make(map[string]*Font)
var selectedFont string

// Load a font created by github.com/ZiRo-/captcha/fontgen
// returns the for usage int AddFont, and an error, if the font can't be loaded.
func LoadFontFromFile(fname string) *Font {
	f := make(map[rune][]byte)

	file, err := os.Open(fname)
	if err != nil {
		return nil
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&f)
	return &Font{f}
}

// Add a font to the internal list of available fonts.
// The name is used to select the font later.
// The first font you add is selected automatically
func AddFont(name string, f *Font) {
	fonts[name] = f
	if len(fonts) == 1 {
		SelectFont(name)
	}
}

// Select a font by the name it was added with AddFont
// If no font with that name exists, the selected font remains unchanged
func SelectFont(name string) {
	_, ok := fonts[name]
	if ok {
		selectedFont = name
	}
}

func Digit2Rune(d byte) rune {
	switch {
	case 0 <= d && d <= 9:
		return rune(d) + '0'
	case 10 <= d && d <= 10+byte('Z'-'A'):
		return rune(d) + 'A' - 10
	case 11+byte('Z'-'A') <= d && d <= 11+byte('Z'-'A')+byte('z'-'a'):
		return rune(d) - 'Z' + 'A' + 'a' - 11
	}
	return 0
}

func Rune2Digit(c rune) byte {
	switch {
	case '0' <= c && c <= '9':
		return byte(c - '0')
	case 'A' <= c && c <= 'Z':
		return byte(c - 'A' + 10)
	case 'a' <= c && c <= 'z':
		return byte(c - 'a' + 'Z' - 'A' + 11)
	}
	return 0
}

func getChar(d byte) []byte {
	if selectedFont == "" {
		panic("No font selected")
	}
	r := Digit2Rune(d)
	return fonts[selectedFont].data[r]
}
