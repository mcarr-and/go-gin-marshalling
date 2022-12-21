## About
Use the [Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework) & [OpenTelemetry](https://opentelemetry.io/docs/) for Observability

The project is a backend service that represents a music store with an in-memory database.

### OpenTelemetry collector 

Send data to the following services:
* Jaeger
* Prometheus

## 1. Running Project

### Docker Compose
[Docker-Compose instructions](docs/Install-Docker-Compose.md)

## K3D
[Local Kubernetes with K3D instructions](docs/Install-Kubernetes-K3D.md)

## TODO
* Run album store externally to Kubernetes with K3D cluster running OpenTelemetry  
* Run album store inside cluster 
* Contract tests for swagger to output 
* Create Prometheus export health endpoint
* Test data builder for creating hundreds of albums for pagination testing and load testing
* Create status endpoint that says if service is up or down.
* Async processing of requests 
* Back pressure on APIs & rate limiting
* pagination of get methods so can receive many and respond in chunks
* Use a database as a data store
* Database migration tooling via scripts 
* Fuzz testing 
* Adding CI server integration
* Adding project to a Docker container
* Helm chart to add Gin Server and Database
* Terrafrom project into EKS

## Bibliography of sites used for creating this project:

Golang tutorial for Gin music store: https://go.dev/doc/tutorial/web-service-gin. 

Go does JSON marshalling and binding in Gin: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

Go Gin testing: https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin

Test benchmarking: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

Gin Examples: https://gin-gonic.com/docs/examples/

Opentelemetry and Gin https://signoz.io/opentelemetry/go/

OpenTelemetry using Otel-collector https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go

OpenTelemetry source of Docker-Compose setup https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo/server

OpenTelemetry unit testing spans https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/test/gintrace_test.go

Go & Docker example https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/Dockerfile