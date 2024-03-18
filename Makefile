.PHONY: build
build:
	go build -v ./app/filmography

.DEFAULT_GOAL : build