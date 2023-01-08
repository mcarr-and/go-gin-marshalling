## About
Use the [Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework) & [OpenTelemetry](https://opentelemetry.io/docs/) for Observability

The project is a backend service that represents a music store with an in-memory database.

### OpenTelemetry collector 

Send data to the following services:
* Jaeger
* Prometheus[TODO]

## 1. Running Project

### Docker Compose

[Docker-Compose fully inclusive instructions](docs/Run-Docker-Compose-Install-Full.md)

[Docker-Compose with album-store as an external application instructions](docs/Run-Docker-Compose-Install-External.md)

## K3D cluster

Run the project with a local Kubernetes cluster with K3D. 

[Local Kubernetes with K3D instructions](docs/Run-K3D.md)

[K3D and using local registry](docs/K3D-registry.md)

[Debugging useful commands](docs/K3D-Debugging.md)

## [WIP & NON-FUNCTIONING] Microk8s Cluster

[Local Kubernetes with Microk8s instructions](docs/Microk8s-Install.md)

## [WIP] Proxy-Service

Standalone server that proxies calls to the album-store

[proxy-service](proxy/README.md)


## TODO
* after front end is in place create sub span from current span to have nested spans. 
* liveness endpoint & wire to album-store
* health endpoint & wire into album-store
* Add prometheus and grafana to cluster
* Write Logs in JSON format
* add all request and response headers and request parameters to the otel attributes.
* Adding CI server integration
* Add open-telmetry to the ingress-nginx so spans are created from the entry point to the cluster.
* prometheus endpoint
* Create status endpoint that says if service is up or down.
* Contract tests compare swagger output to actual output
* Async processing of requests 
* Back pressure on APIs & rate limiting
* Test data builder for creating hundreds of albums for pagination testing and load testing
* pagination of get methods so can receive many and respond in chunks
* Use a database as a data store
* add open-telemetry observability to the database driver.
* Database migration tooling via scripts 
* Fuzz testing 
* Helm chart add Database configuration
* Terraform project into EKS
* Run album store inside cluster via helm chart
* Get working helm chart for album-store - github raw url not working for chart, not in skaffold.
* Run album store externally to Kubernetes with K3D cluster running OpenTelemetry(K3D) GO issue http2 frame too large.

## Bibliography of sites used for creating this project:

Golang tutorial for Gin music store: https://go.dev/doc/tutorial/web-service-gin. 

Go does JSON marshalling and binding in Gin: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

Go Gin testing: https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin

Test benchmarking: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

Gin Examples: https://gin-gonic.com/docs/examples/

Opentelemetry and Gin https://signoz.io/opentelemetry/go/

OpenTelemetry using Protobuf GRPC Otel-collector https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go

OpenTelemetry source of Docker-Compose setup https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo/server

OpenTelemetry unit testing spans https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/test/gintrace_test.go

Go & Docker example https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/Dockerfile

Go & Contract testing for swagger https://github.com/getkin/kin-openapi

Go & HTTP2 + too large frame issue https://kennethjenkins.net/posts/go-nginx-grpc/