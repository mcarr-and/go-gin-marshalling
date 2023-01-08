

This is an example Go-gin application that demonstrates nested spans. 

This uses the opentelemetry instrumented http client[otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp) 

This is a simple pass through service that calls the Album service.

## Prerequisites 
Cluster must have the following deployed
* Jaeger
* Opentelemetry-collector
* Album-Store

# Local Run
## 0. Start Locally
`GRPC_GO_LOG_SEVERITY_LEVEL=info;GRPC_GO_LOG_VERBOSITY_LEVEL=99;INSTANCE_NAME=proxy-service;NAMESPACE=no-namespace;OTEL_LOCATION=localhost:4327;ALBUM_STORE_LOCATION=localhost:9080 go run main.go`

## 1. hit albums url 

[view proxy-service albums](http://localhost:9070/albums)

## 2. view spans

[Jaeger proxy-service spans ](http://localhost:16696/search?service=proxy-service)

## TODO:
* Add to Makefile 
* Add to Docker-compose
* Unit test with mocked out album-store serving different responses.

# Docker Run

## 0. K3D install

[k3D install service ](K3D-Install.md)

## 1. Setup Env

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

add `proxy-service.local` to your list of *.local environments 

`127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local album-store.local proxy-service.local`

