prefix=/usr/local

all:
	go build -buildmode=c-shared -o libgocaptcha.so libgocaptcha.go

install:
	mv libgocaptcha.so $(prefix)/lib/
	mkdir -p $(prefix)/include/gocaptcha
	mv libgocaptcha.h $(prefix)/include/gocaptcha/

test:
	gcc -L$(prefix)/lib -I-L$(prefix)/include -o test.out test.c -lgocaptcha
	
clean:
	rm -f libgocaptcha.so
	rm -f libgocaptcha.h
	rm -f test.out
	rm -f captcha.png

uninstall:
	rm -f $(prefix)/lib/libgocaptcha.so
	rm -rf $(prefix)/include/gocaptcha/
