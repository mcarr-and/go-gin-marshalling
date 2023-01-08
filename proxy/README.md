

This is an example Go-gin application that demonstrates nested spans. 

This uses the opentelemetry instrumented http client [otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp) 

This is a simple pass through service that calls the Album service.

## Prerequisites 
Cluster must have the following deployed
* Jaeger
* Opentelemetry-collector
* Album-Store

## TODO:
* Unit test with mocked out album-store serving different responses.


# Local Run
## 0. Start Locally the Observability service 

```bash
../make docker-compose-limited-start
```

## 1. start album-store

```bash
../make local-start-grpc 
```

## 2. start proxy-service

`GRPC_GO_LOG_SEVERITY_LEVEL=info;GRPC_GO_LOG_VERBOSITY_LEVEL=99;INSTANCE_NAME=proxy-service;NAMESPACE=no-namespace;OTEL_LOCATION=localhost:4327;ALBUM_STORE_LOCATION=localhost:9080 go run main.go`

## 3. hit albums url 

[view proxy-service albums](http://localhost:9070/albums)

## 4. view spans

[Jaeger proxy-service spans ](http://localhost:16696/search?service=proxy-service)

# Docker Run

## 1. K3D install

[k3D install service ](K3D-Install.md)

## 2. Setup Env

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

add `proxy-service.local` to your list of *.local environments 

`127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local album-store.local proxy-service.local`

## 3. install album-store 

```bash
../make k3d-album-deploy-deployment;
```

## 4. install proxy-service

```bash
make k3d-proxy-deploy-deployment;
```

## 5. View the events in the different Services in K3D

[View Jaeger](http://jaeger.local:8070/search?limit=20&service=album-store)

[View Kubernetes environment](http://k-dashboard:8070/)

## 6. Run Some Tests

### 6.1 curl

```bash
curl --insecure --location 'http://proxy-service.local:8070/albums/'; 
```

### 6.3 Run Test Suite

```bash
make k3d-test;
```

### 6.2 Postman

[Postman files](../test/postman_collection.json)

1. Import the folder `../test`
1. Set Environment to `proxy-service.local`
1. Open a test in the `Album-Store` collection and run it.

## 7. Stop album-store server & Services

Ctr + C on the terminal window where you started `make skaffold-dev`

## 8. Delete Proxy-service

```bash
make k3d-proxy-undeploy-deployment;
```

## 9. Delete Album-Store

```bash
../make k3d-album-undeploy-deployment;
```

## 9. Delete Cluster

```bash
../make k3d-cluster-delete
```