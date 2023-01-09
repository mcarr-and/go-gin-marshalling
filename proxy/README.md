

This is an example Go-gin application that demonstrates nested spans. 

This uses the opentelemetry instrumented http client [otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp) 

This is a simple pass through service that calls the Album service.

## Prerequisites 
Cluster must have the following deployed
* Jaeger
* Opentelemetry-collector
* Album-Store

## TODO:
* FLush out unit tests for examples.
* need to improve how the mock interacts with the code under test to have correct number of spans. 
  * 1 successful call has 3 spans.  

# Run 

## Docker-Compose

[Docker-Compose](Run-Docker-Compose-Limited.md)

# K3d Run

[k3D install service ](Run-K3D-Install.md)
