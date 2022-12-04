Test project to understand the  [Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework)

Project is a backend service that music store with an in-memory database.

Using Jaeger tracing for observability https://www.jaegertracing.io/


## 0. Expected tooling to run this project

1. Go
2. Docker 

## 1. Start All Expected Services 

The following will start 
* OpenTelemetry collector
* Jaeger
* Zipkin
* Prometheus 

 
```bash
   make docker-compose-start 
```

## 2. Start Server

```bash
  make start;
```

## 3. Run Tests


### (Either) PostMan collection for tests

1.Import the Postman collection into your Postman to run the tests. 

[Postman Collection](test/Album-Store.postman_collection.json)

2. Run tests

Postman Collection
   
### (Or) Run a few Curl commands 

```bash
  make local-run-tests;
```

## 4. View the Data in the systems

[Jaeger](http://localhost:16686/search?limit=20&service=album-store)

[Zipkin](http://localhost:9411/zipkin/)

[Prometheus](http://localhost:9090/graph?g0.expr=%7Bjob%3D~%22.%2B%22%7D%20&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h)

[View Albums](http://localhost:9080/albums)

## 5. Stop server & Jaeger 

1. Stop Server

`Ctr + C` in the terminal window where go is running. 

2. Stop Jaeger

```bash
  make docker-comose-stop;
```

## Project includes:

* Opentelemetry via OpenTelemetry Collector.
  
* Testing
  * Unit testing sending and receiving JSON to Gin
    * Get all albums
    * Get album by ID
    * Get album by ID that is not found
    * Post to create new album
    * Post album without all the required JSON fields to be a valid object to Gin & V10
  * Benchmark tests for throughput for all unit tests
* Docker build
* Docker-Compose to use local tooling  
  * Jaeger
  * Zipkin 
  * Prometheus 

## TODO
* Prometheus export health endpoint
* Swagger
* Use a database as a data store
* Database migration tooling via scripts
* Fuzz testing 
* Unit test the OpenTelemetry messages are received
* Adding CI server integration
* Adding project to a Docker container
* K3D cluster to run Docker container in K8s.
* Helm chart to add Gin Server and Database
* Skaffold to set up the K8s cluster for this project.

## Source Sites for creating this project

Golang tutorial for Gin music store: https://go.dev/doc/tutorial/web-service-gin. 

Go does JSON marshalling and binding in Gin: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

Go Gin testing: https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin

Test benchmarking: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

Gin Examples: https://gin-gonic.com/docs/examples/

Opentelemetry and Gin https://signoz.io/opentelemetry/go/

OpenTelemetry using Otel-collector https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go

OpenTelemetry source of Docker-Compose setup https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo/server