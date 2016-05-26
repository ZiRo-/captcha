// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// capgen is an utility to test captcha generation.
package main

import (
	"flag"
	"fmt"
	"github.com/ZiRo-/captcha"
	"github.com/ZiRo-/captcha/libgocaptcha"
	"io"
	"log"
	"os"
)

var (
	flagLen  = flag.Int("len", captcha.DefaultLen, "length of captcha")
	flagImgW = flag.Int("width", libgocaptcha.StdWidth, "image captcha width")
	flagImgH = flag.Int("height", libgocaptcha.StdHeight, "image captcha height")
	fontFile = flag.String("ff", "Monospace.gob", "font file")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: captcha [flags] filename\n")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	fname := flag.Arg(0)
	if fname == "" {
		usage()
		os.Exit(1)
	}

	fn := libgocaptcha.LoadFontFromFile(*fontFile)
	if fn == nil {
		log.Fatalf("Couldn't load font file")
	}
	libgocaptcha.AddFont("font", fn)

	f, err := os.Create(fname)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer f.Close()
	var w io.WriterTo
	d := libgocaptcha.RandomDigits(*flagLen)
	w = libgocaptcha.NewImage("", d, *flagImgW, *flagImgH)
	_, err = w.WriteTo(f)
	if err != nil {
		log.Fatalf("%s", err)
	}
	for _, c := range d {
		fmt.Printf("%c", libgocaptcha.Digit2Rune(c))
	}
	fmt.Println()
}
