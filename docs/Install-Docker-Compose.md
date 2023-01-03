## 0. Expected tooling to run this project

1. Go
2. Docker 

## 1. Start All Observability & Log Viewing Services
 
```bash
make docker-compose-start;
```

## 2. Start album-store Go/Gin Server with flags set

```bash
	export NAMESPACE=no-namespace;
	export INSTANCE_NAME=album-store-1;
	export OTEL_LOCATION=localhost:4327;
```

```bash
make local-start-docker-compose-grpc;
```

#### Note: the application will not start without the OpenTelemetry collector running

## 3. Run Some Tests

[Postman Collection](../test/Album-Store.postman_collection.json)

```bash
make local-test;
```

## 4. View the events in the different Services

[View Jaeger](http://localhost:16696/search?limit=20&service=album-store)

[View Prometheus](http://localhost:9090/graph?g0.expr=%7Bjob%3D~%22.%2B%22%7D%20&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h)

## 5. Stop album-store server & Services  

### 1. Stop Server

`Ctr + C` in the terminal window where go is running. 

### 2. Stop Observability and Log Viewing Services

```bash
make docker-compose-stop;
```