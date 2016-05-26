// Copyright 2016 ZiRo. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"github.com/ZiRo-/captcha/libgocaptcha"
	"image/png"
	"unsafe"
)

/*
extern int libgocaptcha_load_font(char* fname, char *fontname);
extern void libgocaptcha_select_font(char* name);
extern void libgocaptcha_random_digits(unsigned char* out, int len);
extern void libgocaptcha_set_char_range(unsigned char b);
extern unsigned char libgocaptcha_rune2digit(unsigned int r);
extern unsigned int libgocaptcha_digit2rune(unsigned char d);
extern int libgocaptcha_new_image(char* id, unsigned char* digits, int dlen, int width, int height, unsigned char* buffer, int len);
*/
import "C"


//export libgocaptcha_load_font
func libgocaptcha_load_font(fname *C.char, fontname *C.char) C.int {
	go_fname := C.GoString(fname)
	go_fontname := C.GoString(fontname)
	fn := libgocaptcha.LoadFontFromFile(go_fname)
	if fn == nil {
		return 1
	}
	libgocaptcha.AddFont(go_fontname, fn)
	return 0
}

//export libgocaptcha_select_font
func libgocaptcha_select_font(name *C.char) {
	go_name := C.GoString(name)
	libgocaptcha.SelectFont(go_name)
}

//export libgocaptcha_random_digits
func libgocaptcha_random_digits(out *C.uchar, length C.int) {
	bs := libgocaptcha.RandomDigits(int(length))
	go_out := (*[1 << 30]C.uchar)(unsafe.Pointer(out))[:length:length]
	for i, b := range(bs) {
		go_out[i] = C.uchar(b)
	}
}

//export libgocaptcha_set_char_range
func libgocaptcha_set_char_range(r C.uchar) {
	go_r := byte(r)
	libgocaptcha.SetCharacterRange(go_r)
}

//export libgocaptcha_rune2digit
func libgocaptcha_rune2digit(r C.uint) C.uchar {
	return C.uchar(libgocaptcha.Rune2Digit(rune(r)))
}

//export libgocaptcha_digit2rune
func libgocaptcha_digit2rune(r C.uchar) C.uint {
	return C.uint(libgocaptcha.Digit2Rune(byte(r)))
}

//export libgocaptcha_new_image
func libgocaptcha_new_image(id *C.char, digits *C.uchar, dlength C.int, width, height C.int, buffer *C.uchar, length C.int) C.int {
	go_id := C.GoString(id)
	go_digits := C.GoBytes(unsafe.Pointer(digits), dlength) 
	go_img := libgocaptcha.NewImage(go_id, go_digits, int(width), int(height))
	
	go_buffer := (*[1 << 30]C.uchar)(unsafe.Pointer(buffer))[:length:length]
	
	
	var buf bytes.Buffer
	if err := png.Encode(&buf, go_img); err != nil {
		return -2
	}
	
	bs := buf.Bytes()
	l := len(bs)
	if l > int(length) {
		return -1
	}
	
	for i, b := range(bs) {
		go_buffer[i] = C.uchar(b)
	}
	return C.int(l)
}


func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C shared library
}
