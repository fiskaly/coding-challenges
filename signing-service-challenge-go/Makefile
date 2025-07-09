# 
# Small Makefile for the test project
#
# Raffaele Russo
#

.PHONY: build run check 

build:
	go mod tidy && \
	go build

run:
	go run main.go

check:
	go test ./... -v -p 1 -coverage