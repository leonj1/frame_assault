.PHONY: build run

build:
	go build -o frame_assault main.go

run: build
	./frame_assault