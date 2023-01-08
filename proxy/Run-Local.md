# Local Run
## 0. Start Locally the Observability service

```bash
  ../make docker-compose-limited-start;
```

## 1. start album-store

```bash
  ../make local-start-grpc; 
```

## 2. start proxy-service

`GRPC_GO_LOG_SEVERITY_LEVEL=info;GRPC_GO_LOG_VERBOSITY_LEVEL=99;INSTANCE_NAME=proxy-service;NAMESPACE=no-namespace;OTEL_LOCATION=localhost:4327;ALBUM_STORE_LOCATION=localhost:9080 go run main.go`

## 3. Run Tests

### 3.1 Command line 

```bash
curl --insecure --location 'http://localhost:9070/albums/'; 
```

```bash
  make local-test;
```

### 3.2 via browser

[view proxy-service albums](http://localhost:9070/albums)

### 3.3 via Postman 

[Postman files](../test/postman_collection.json)

1. Import the folder `../test`
1. Set Environment to `localhost:9070`
1. Open a test in the `Album-Store` collection and run it.


## 4. view spans

These will have 2 spans.

* 1 for Proxy-Service
* 1 for Album-Store

[Jaeger proxy-service spans](http://localhost:16696/search?service=proxy-service)
