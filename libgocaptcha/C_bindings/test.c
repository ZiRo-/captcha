#include <gocaptcha/libgocaptcha.h>
#include <stdio.h>
#include <stdlib.h>

int main() {
	int err = libgocaptcha_load_font("../../fontgen/UbuntuMono/example_font_UbuntuMono.gob", "Ubuntu Mono");
	
	if (err) {
		printf("Couldn't load font\n");
		return 1;
	}
	
	unsigned char digits[6];
	libgocaptcha_random_digits(digits, 6);
	
	
	int size = 10 * 1024 * 1024; // 10 kB
	
	unsigned char *buffer = (unsigned char*) malloc(sizeof(unsigned char)*size);
	
	int len = libgocaptcha_new_image("1", digits, 6, 350, 150, buffer, size);
	
	if (len<0) {
		printf("Buffer too small\n");
		return 1;
	}
	
	FILE* f = fopen("captcha.png","wb");

	if (f){
		fwrite(buffer, len, 1, f);
	}
	else{
		printf("Couldn't create file\n");
		return 1;
	}
	
	free(buffer);
	fclose(f);
}
