Test project to understand the  [Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework)

Project is a backend service that music store with an in-memory database.

Using Jaeger tracing for observability https://www.jaegertracing.io/


## 0. Expected tooling to run this project

1. Go
2. Docker 

## 1. Start Jaeger and server

1. Install and Run Docker Jaeger
```bash
   make jager-install; 
```

[View Albums](http://localhost:9080/albums)

[Jaeger tracing](http://localhost:16686/search )

2. Start Server

```bash
  make start;
```
## 2. Run Tests


### (Either) PostMan collection for tests

1.Import the Postman collection into your Postman to run the tests. 

[Postman Collection](test/Album-Store.postman_collection.json)

2. Run tests

Postman Collection
   
### (Or) Run a few Curl commands 

`curl http://localhost:9080/albums/1`

`curl http://localhost:9080/albums/666`

`curl http://localhost:9080/albums`


## 3. See the Spans for the album-store

[View Jaeger spans for the tests that you have run album-store](http://localhost:16686/search?limit=20&lookback=1h&maxDuration&minDuration&service=album-store)

## 4. Stop server & Jaeger 

1. Stop Server

`Ctr + C` in the terminal window where go is running. 

2. Stop Jaeger

```bash
  make jaeger-stop;
```

## Project includes:

* WIP
  * Opentelemetry via Jaeger
  
* Unit testing sending and receiving JSON to Gin
  * Get all albums
  * Get album by ID
  * Get album by ID that is not found
  * Post to create new album
  * Post album without all the required JSON fields to be a valid object to Gin & V10
* Benchmark tests for throughput for all unit tests
* Docker

## TODO
* Prometheus export data collection
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