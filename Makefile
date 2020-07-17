all: build generate

setup: setup.sh
	./setup.sh

build: src
	cd src && go build -o ../bsb

generate:
	./bsb -target-dir texts -min 1000 > output.txt
