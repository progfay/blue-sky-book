all: build generate

setup: setup.sh
	./setup.sh

build: src/main.go
	go build -o bsb ./src/main.go

generate:
	./bsb > output.txt
