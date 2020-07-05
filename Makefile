all: build generate

setup: setup.sh
	./setup.sh

build: main.go
	go build -o bsb main.go

generate:
	./bsb > output.txt
