.PHONY: build

build:
	go build -v -o go-gin-example

test:
	go test -v

test-benchmark:
	go test -bench=. -count 2 -run=^# -benchmem

local-start-example:
	go mod tidy
	go run main.go;

docker-build:
	docker build -t go-gin-example:0.1 .

docker-start-example:
	docker run -d -p 9080:9080 --name go-gin-example go-gin-example:0.1

docker-stop-example:
	docker stop go-gin-example

docker-compose-start:
	docker-compose up -d

docker-compose-stop:
	docker-compose down

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html