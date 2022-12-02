.PHONY: build

build:
	go build -v -o go-gin-example

test:
	go test -v

benchmark:
	go test -bench=. -count 2 -run=^# -benchmem

jaeger-install:
	 docker run -d --name jaeger --label jaeger-all-in-one \
          -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
          -e COLLECTOR_OTLP_ENABLED=true \
          -p 5775:5775/udp \
          -p 6831:6831/udp \
          -p 6832:6832/udp \
          -p 5778:5778 \
          -p 16686:16686 \
          -p 14268:14268 \
          -p 9411:9411 \
          -p 4317:4317 \
          -p 4318:4318 \
          jaegertracing/all-in-one:1.39;

jaeger-start:
	docker start $(shell docker ps -qa -f "name=jaeger" -f "label=jaeger-all-in-one" );

jaeger-stop:
	docker stop $(shell docker ps -qa -f "name=jaeger" -f "label=jaeger-all-in-one" );

start:
	go mod tidy
	go run main.go;

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html