
.PHONY: all clean prepare test build install uninstall

all: clean build

clean:
	rm -rf bin/*

prepare:
	go mod tidy
	go fmt .
	go vet .

test:
	go test ./...

build: prepare test
	go build -o bin/wol .

install: clean build
	sudo cp bin/wol /usr/local/bin/

uninstall:
	sudo rm -f /usr/local/bin/wol

