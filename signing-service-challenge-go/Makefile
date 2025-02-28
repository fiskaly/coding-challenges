.PHONY: all build test run clean

all: test build

build:
	go build -o signing-service-challenge .

test:
	go test ./... -v -race

run: build
	./signing-service-challenge

clean:
	rm -f signing-service-challenge
