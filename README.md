Package captcha
=====================
[![GoDoc](https://godoc.org/github.com/ZiRo-/captcha?status.svg)](https://godoc.org/github.com/ZiRo-/captcha)

This is a fork of https://github.com/dchest/captcha
Original by Dmitry Chestnykh

	import "github.com/ZiRo/captcha"

Package captcha implements generation and verification of image CAPTCHAs.

A captcha solution is the sequence of digits 0-9 and letters a-z, A-Z 
with a defined length.

The captcha is a PNG-encoded image with the solution printed on
it in such a way that makes it hard for computers to solve it using OCR.

This package only requires font files. See [github.com/ZiRo-/captcha/fontgen](https://github.com/ZiRo-/captcha/tree/master/fontgen)
for details on how to get them.
So, before you start generating captchas, you have to load a font:
``` go
	fn := libgocaptcha.LoadFontFromFile("Ubuntu-Mono.gob)
	if fn == nil {
		log.Fatalf("Couldn't load font file")
	}
	libgocaptcha.AddFont("font", fn)
	libgocaptcha.SetCharacterRange(libgocaptcha.MODULE_UPPER)
```

To make captchas one-time, the package includes a memory storage that stores
captcha ids, their solutions, and expiration time. Used captchas are removed
from the store immediately after calling Verify or VerifyString, while
unused captchas (user loaded a page with captcha, but didn't submit the
form) are collected automatically after the predefined expiration time.
Developers can also provide custom store (for example, which saves captcha
ids and solutions in database) by implementing Store interface and
registering the object with SetCustomStore.

Captchas are created by calling New, which returns the captcha id. Their
representations, though, are created on-the-fly by calling WriteImage. 
Created representations are not stored anywhere, but
subsequent calls to these functions with the same id will write the same
captcha solution. Reload function will create a new different solution for
the provided captcha, allowing users to "reload" captcha if they can't solve
the displayed one without reloading the whole page.  Verify and VerifyString
are used to verify that the given solution is the right one for the given
captcha id.

Server provides an http.Handler which can serve image and audio
representations of captchas automatically from the URL. It can also be used
to reload captchas.  Refer to Server function documentation for details, or
take a look at the example in [github.com/ZiRo-/captcha/capexample](https://github.com/ZiRo-/captcha/tree/master/capexample)

C/C++ bindings
--------

See https://github.com/ZiRo-/captcha/tree/master/libgocaptcha/C_bindings

Examples
--------

![Image](https://github.com/ZiRo-/captcha/raw/master/capgen/example.png)

