// Copyright 2011-2014 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package libgocaptcha

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// These constants are used to set the used character ranges in the captcha image
const (
	MODULE_DIGIT = 10
	MODULE_UPPER = 36
	MODULE_LOWER = 62
)

var rand_mod byte = MODULE_UPPER

// idLen is a length of captcha id string.
// (20 bytes of 62-letter alphabet give ~119 bits.)
const idLen = 20

// idChars are characters allowed in captcha id.
var idChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// rngKey is a secret key used to deterministically derive seeds for
// PRNGs used in image and audio. Generated once during initialization.
var rngKey [32]byte

func init() {
	if _, err := io.ReadFull(rand.Reader, rngKey[:]); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
}

// Purposes for seed derivation. The goal is to make deterministic PRNG produce
// different outputs for images and audio by using different derived seeds.
const (
	imageSeedPurpose = 0x01
	audioSeedPurpose = 0x02
)

// Set the desired character ranges for the captcha image. For example:
//	captcha.SetCharacterRange(captcha.MODULE_LOWER)
// would use 0-9, A-Z and a-z
// The default is MODULE_UPPER which means 0-9 and A-Z
func SetCharacterRange(rang byte) {
	if rang == MODULE_DIGIT || rang == MODULE_UPPER || rang == MODULE_LOWER {
		rand_mod = rang
	}
}

// deriveSeed returns a 16-byte PRNG seed from rngKey, purpose, id and digits.
// Same purpose, id and digits will result in the same derived seed for this
// instance of running application.
//
//   out = HMAC(rngKey, purpose || id || 0x00 || digits)  (cut to 16 bytes)
//
func deriveSeed(purpose byte, id string, digits []byte) (out [16]byte) {
	var buf [sha256.Size]byte
	h := hmac.New(sha256.New, rngKey[:])
	h.Write([]byte{purpose})
	io.WriteString(h, id)
	h.Write([]byte{0})
	h.Write(digits)
	sum := h.Sum(buf[:0])
	copy(out[:], sum)
	return
}

// RandomDigits returns a byte slice of the given length containing
// pseudorandom numbers in range 0-module. The slice can be used as a captcha
// solution.
func RandomDigits(length int) []byte {
	return randomBytesMod(length, rand_mod)
}

// randomBytes returns a byte slice of the given length read from CSPRNG.
func randomBytes(length int) (b []byte) {
	b = make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return
}

// randomBytesMod returns a byte slice of the given length, where each byte is
// a random number modulo mod.
func randomBytesMod(length int, mod byte) (b []byte) {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb := 255 - byte(256%int(mod))
	b = make([]byte, length)
	i := 0
	for {
		r := randomBytes(length + (length / 4))
		for _, c := range r {
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = c % mod
			i++
			if i == length {
				return
			}
		}
	}

}

// RandomId returns a new random id string.
func RandomId() string {
	b := randomBytesMod(idLen, byte(len(idChars)))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b)
}
