.PHONY: build
build:
	go build -v -o go-gin-example

test:
	go test -v

benchmark:
	go test -bench=. -count 2 -run=^# -benchmem

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html