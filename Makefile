
ALL_SRC = cmd/welp/welp.go *.go

all: welp test

run: welp
	./welp

welp: $(ALL_SRC)
	go build ./cmd/welp/welp.go

test:
	go test *.go

ctest:
	go test -cover *.go
	go tool cover -func=coverage.out

wtest:
	go test -cover *.go
	go tool cover -html=coverage.out
