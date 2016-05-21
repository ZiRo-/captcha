// Copyright 2016 ZiRo. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"encoding/gob"
	"os"
)

const (
	fontWidth  = 15
	fontHeight = 25
	blackChar  = 1
)

type Font map[rune][]byte

var fonts = make(map[string]Font)
var selectedFont string

func LoadFontFromFile(fname string) (Font, error) {
	f := make(map[rune][]byte)

	file, err := os.Open(fname)
	if err != nil {
		return f, err
	}
	dec := gob.NewDecoder(file)
	err = dec.Decode(&f)
	return f, err
}

func AddFont(name string, f Font) {
	fonts[name] = f
	if len(fonts) == 1 {
		SelectFont(name)
	}
}

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
	case 10 <= d && d <= 10+byte('z'-'a'):
		return rune(d) + 'a' - 10
	case 11+byte('z'-'a') <= d && d <= 11+byte('z'-'a')+byte('Z'-'A'):
		return rune(d) - 'z' + 'a' + 'A' - 11
	}
	return 0
}

func Rune2Digit(c rune) byte {
	switch {
	case '0' <= c && c <= '9':
		return byte(c - '0')
	case 'a' <= c && c <= 'z':
		return byte(c - 'a' + 10)
	case 'A' <= c && c <= 'Z':
		return byte(c - 'A' + 'z' - 'a' + 11)
	}
	return 0
}

func getChar(d byte) []byte {
	if selectedFont == "" {
		panic("No font selected")
	}
	r := Digit2Rune(d)
	return fonts[selectedFont][r]
}
