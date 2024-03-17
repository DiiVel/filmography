.PHONY: build
build:
	go build -v ./cmd/filmography

.DEFAULT_GOAL : build