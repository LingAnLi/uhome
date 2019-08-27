
GOPATH:=$(shell go env GOPATH)

.PHONY: proto test docker


build:

	go build -o UHomeWeb-web main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t UHomeWeb-web:latest
