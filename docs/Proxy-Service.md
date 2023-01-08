

This is an example Go-gin application that demonstrates nested spans. 

This uses the opentelemetry instrumented http client[otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp) 

This is a simple pass through service that calls the Album service.


## 0. Start Locally
`GRPC_GO_LOG_SEVERITY_LEVEL=info;GRPC_GO_LOG_VERBOSITY_LEVEL=99;INSTANCE_NAME=proxy-service;NAMESPACE=no-namespace;OTEL_LOCATION=localhost:4327;ALBUM_STORE_URL=http://localhost:9080 go run main.go`

## 1. hit albums url 

[view proxy-service albums](http://localhost:9070/albums)

## 2. view spans

[Jaeger proxy-service spans ](http://localhost:16696/search?service=proxy-service)

## TODO:
* Add to Makefile 
* Add to Docker-compose
* Unit test with mocked out album-store serving different responses.