# Local with Docker-Compose

## 0. start Docker-Compose

Terminal window 1:

```bash
  ../make docker-compose-limited-start
```

## 1. Start Album-store

Terminal window 2:

```bash
../make local-start-grpc;
```

## 2. Start Proxy-Service

Terminal window 3:

```bash
make local-start-grpc;
```

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

## 4. view spans in Jaeger

These will have 2 spans.

* 1 for Proxy-Service
* 1 for Album-Store

[Jaeger proxy-service spans](http://localhost:16696/search?service=proxy-service)
