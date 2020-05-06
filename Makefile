all: build

.PHONY: build
build: build_games
	go build

.PHONY: run
run: build_games
	go run main.go

.PHONE: test
test: build_games
	go test -v ./...

.PHONY: build_games
build_games:
	cd games; python gen.py; gofmt -w games.go
