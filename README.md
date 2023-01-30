## About
Use the [Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework) & [OpenTelemetry](https://opentelemetry.io/docs/) for Observability

The project is a backend service that represents a music store with an in-memory database.

### OpenTelemetry Collector 

Sends data to the following services:
* Jaeger
* Prometheus

## Proxy-Service

Standalone server that proxies calls to the `album-store`

Used for showing nested spans in open-telemetry.

Proxy-Service uses the `otelhttp.client` to make requests which produces nested spans

[proxy-service](proxy/.)


## TL;DR
Run the following, so you can see how the services work and produce nested OpenTelemetry spans.

```bash 
  make && make local-proxy-test;
```

http://localhost:16696/ <- Jaeger to see spans. 

Pick the `proxy-service` from the `service` dropdown to see nested spans. 


## Running Project

### Docker Compose

[Docker-Compose fully inclusive instructions](docs/Run-Docker-Compose-Install-Full.md)

[Docker-Compose with album-store as an external application instructions](docs/Run-Docker-Compose-Install-Limited.md)

### K3D cluster

Run the project with a local Kubernetes cluster with K3D. 

[Local Kubernetes with K3D instructions](docs/Run-K3D.md)

[K3D and using local registry](docs/K3D-registry.md)

[Debugging useful commands](docs/K3D-Debugging.md)

### [WIP & NON-FUNCTIONING] Microk8s Cluster

[Local Kubernetes with Microk8s instructions](docs/Microk8s-Install.md)


## TODO 
* Grafana dashboards to display the GO services, Jaeger, nginx-ingress
* Write Logs in JSON format
* Adding CI server integration
* Add open-telemetry to the ingress-nginx so spans are created from the entry point to the cluster.
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
* Terraform project into EKS or GKE
* Run album store externally to Kubernetes with K3D cluster running OpenTelemetry(K3D) GO issue http2 frame too large.

## Bibliography of sites used for creating this project:

Golang tutorial for Gin music store: https://go.dev/doc/tutorial/web-service-gin. 

Go does JSON marshalling and binding in Gin: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

Go Gin testing: https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin

Otel gin testing https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/test/gintrace_test.go

Test benchmarking: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

Gin Examples: https://gin-gonic.com/docs/examples/

Opentelemetry and Gin https://signoz.io/opentelemetry/go/

OpenTelemetry using Protobuf GRPC Otel-collector https://github.com/open-telemetry/opentelemetry-go/blob/main/example/otel-collector/main.go

OpenTelemetry source of Docker-Compose setup https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo/server

OpenTelemetry unit testing spans https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/test/gintrace_test.go

Go & Docker example https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin/example/Dockerfile

Go & Contract testing for swagger https://github.com/getkin/kin-openapi

Go & HTTP2 + too large frame issue https://kennethjenkins.net/posts/go-nginx-grpc/

Full OpenTelemetry Demo in multiple languages https://github.com/open-telemetry/opentelemetry-demo

Grafana import dashboards via configmaps https://blog.cloudcover.ch/posts/grafana-helm-dashboard-import/

[Grafana dashboards imported from configmaps](./docs/Grafana-Helm-Prometheus-Configmap-Dashboards.md)