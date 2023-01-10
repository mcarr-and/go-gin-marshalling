# Local with Docker-Compose

## 0. start Docker-Compose

Terminal window 1:

```bash
  make docker-compose-limited-start
```

## 1. Start Album-store

Terminal window 2:

```bash 
  make local-start-album-grpc;
```

## 2. Start Proxy-Service


### 2.1 Start via command line
Terminal window 3:

```bash
  make local-start-proxy-grpc;
```

### Start via IDE 

`GRPC_GO_LOG_SEVERITY_LEVEL=info;GRPC_GO_LOG_VERBOSITY_LEVEL=99;INSTANCE_NAME=proxy-service;NAMESPACE=no-namespace;OTEL_LOCATION=localhost:4327;ALBUM_STORE_URL=http://localhost:9080 go run main.go`


## 3. Run tests

### 3.1 Command line 
```bash
  curl --insecure --location 'http://localhost:9070/albums/'; 
```

```bash
  make local-test;
```

### 3.2 Browser 

[view proxy-service albums](http://localhost:9070/albums)

### 3.3 via Postman

[Postman files](../test/postman_collection.json)

1. Import the folder `../test`
1. Set Environment to `localhost:9070`
1. Open a test in the `Album-Store` collection and run it.


## 4. view spans in Jaeger

These will have 2 spans.

* 1 for Proxy-Service
* 1 for Album-Store

[Jaeger proxy-service spans](http://localhost:16696/search?service=proxy-service)
