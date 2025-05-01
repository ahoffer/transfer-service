.PHONY: clean build test docker-build

clean:
	go clean

build:
	go build ./...

test:
	go test ./...

docker-build: test
	docker build -t transfer-service .

all: clean build test docker-build 