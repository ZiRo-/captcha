#Go Captcha - C bindings

## Installation
Grab Go from https://golang.org/dl/ or your package manager, and set it up: https://golang.org/doc/install#testing
Once that's done just run
```
$ go get github.com/ZiRo-/captcha/libgocaptcha
$ git clone https://github.com/ZiRo-/captcha
$ cd captcha/libgocaptcha/C_bindings
$ make
$ sudo make install
```

## Usage
Check test.c on how to use it.

The availbe functions are:
```
extern int libgocaptcha_load_font(char* fname, char *fontname);
extern void libgocaptcha_select_font(char* name);
extern void libgocaptcha_random_digits(unsigned char* out, int len);
extern void libgocaptcha_set_char_range(unsigned char b);
extern unsigned char libgocaptcha_rune2digit(unsigned int r);
extern unsigned int libgocaptcha_digit2rune(unsigned char d);
extern int libgocaptcha_new_image(char* id, unsigned char* digits, int dlen, int width, int height, unsigned char* buffer, int len);
```
