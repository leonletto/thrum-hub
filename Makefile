.PHONY: build test clean install

BINARY := bin/thrumhub
PKG := github.com/leonletto/thrum-hub

build:
	go build -o $(BINARY) ./cmd/thrumhub

test:
	go test ./...

clean:
	rm -f $(BINARY)

install: build
	cp $(BINARY) $(HOME)/.local/bin/thrumhub
