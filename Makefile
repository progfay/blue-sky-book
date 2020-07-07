all: build generate

setup: setup.sh
	./setup.sh

build: src
	cd src && go build -o ../bsb

generate:
	./bsb -target-dir texts -min 50 -max 80 > output.txt
