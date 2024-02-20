clean:
	rm logs/*
	rm WebAuthnSSO

build:
	go build .

dev: build
	./WebAuthnSSO -debug

run: build
	./WebAuthnSSO