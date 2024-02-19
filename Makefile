.PHONY: build
APP_NAME=gd


build:
	go mod tidy
	mkdir -p bin
	go build -o bin/$(APP_NAME) main.go
	chmod +x bin/$(APP_NAME)

