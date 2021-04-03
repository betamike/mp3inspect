.PHONY: test build

test/files:
	cd test && ./generate_mp3s.sh && ./generate_golden_files.sh
test: test/files
	go test ./...
build:
	go build -o bin/mp3inspect inspector.go

.DEFAULT_GOAL := build
