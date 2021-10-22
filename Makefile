
.PHONY: all

all: build run

.PHONY: build
.PHONY: run

build:
	@go build -o ./bin/huer ./cmd/hue

run:
	@./bin/huer