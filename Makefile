.PHONY: build

build:
	go build -ldflags "-X main.version=0.1 -X main.gitHash=`git rev-parse --short HEAD`" -v -o go-gin-example

test:
	go test -v

test-benchmark:
	go test -bench=. -count 2 -run=^# -benchmem

local-start: build
	./go-gin-example -otel-location=localhost:4327 -namespace=no-namespace -instance-name=go-gin-example-1

local-test:
	curl --location --request GET 'http://localhost:9080/albums/1' --header 'Accept: application/json';
	curl --location --request GET 'http://localhost:9080/albums/666' --header 'Accept: application/json';
	curl --location --request GET 'http://localhost:9080/albums/X' --header 'Accept: application/json';
	curl --location --request GET 'http://localhost:9080/albums';
	curl --location --request POST 'http://localhost:9080/albums' \
		--header 'Content-Type: application/json' --header 'Accept: application/json' \
		--data-raw '{ "XID": 10, "Titlexx": "Blue Train", "Artistx": "John Coltrane", "Price": 56.99, "X": "asdf" }';
	curl --location --request POST 'http://localhost:9080/albums' \
    		--header 'Content-Type: application/json' --header 'Accept: application/json' \
    		--data-raw '{ "ID": -1, "Title": "s", "Artist": "p", "Price": -0.1}';
	curl --location --request POST 'http://localhost:9080/albums' \
        --header 'Content-Type: application/json' --header 'Accept: application/json' \
        --data-raw '{"ID": 10, "Title": "The Ozzman Cometh", "Artist": "Black Sabbath", "Price": 66.60}';

docker-build:
	docker build -t go-gin-example:0.1 .

docker-start:
	docker run -d -p 9080:9080 --name go-gin-example go-gin-example:0.1

docker-stop:
	docker stop go-gin-example

docker-compose-start:
	docker-compose up -d

docker-compose-stop:
	docker-compose down

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html