#!/bin/bash

go build -buildmode=c-archive -o libgocaptcha.a libgocaptcha.go
gcc -o test test.c libgocaptcha.a -lpthread
