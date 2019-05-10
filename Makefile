
ALL_SRC = cmd/welp/welp.go lexer/*.go parser/*.go *.go

all: welp test

run: welp
	./welp

welp: $(ALL_SRC)
	go build ./cmd/welp/welp.go

test:
	go test ./...

ctest:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

wtest:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
