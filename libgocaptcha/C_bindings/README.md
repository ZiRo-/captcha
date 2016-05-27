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

## Performance

C implmentation (test.c)

```
N = 10000
Avg. time:  0.035416
Stddev:     0.003913
Max. time:  0.201672
Min. time:  0.026584
```


Go implementation (capgen)

```
N = 10000
Avg. time:  0.035975
Stddev:     0.004922
Max. time:  0.244690
Min. time:  0.026796
```

