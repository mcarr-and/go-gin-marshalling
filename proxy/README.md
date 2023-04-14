

This is an example Go-gin application that demonstrates nested spans. 

This uses the opentelemetry instrumented http client [otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/net/http/otelhttp) 

This is a simple pass through service that calls the Album service.

## Prerequisites 
Cluster must have the following deployed
* Jaeger
* Opentelemetry-collector
* Album-Store

## TODO:
* version file for go?

# Run 

## Docker-Compose

[Docker-Compose](../docs/Run-Docker-Compose-Install-Limited.md)

# K3d Run

[k3D install service ](../docs/K3D-run.md)
