## Ingress traffic

For infrastructure to be available externally from the cluster you need to add 3 definitions per service.

VirtualService = Maps the Service in kubernetes into the Isto.

DestinationRule = how do you get there

Gateway = Expose the VirtualService in Istio to the outside


## Exposing infrastructure for this project

This is created in the [helm-repo istio-ingress-charts](../install/helm/istio-ingress-charts) for the following infrastructure
* [Kiali](http://kiali.local:8070)
* [Prometheus](http://prometheus.local:8070)
* [Grafana](http://grafana.local:8070)
* [Jaeger all in one](http://jaeger.local:8070)
* [Kubernetes Dashboard](http://k-dashboard.local:8070)

## Traffic inside cluster = inside mesh

Mesh service = VirtualService(can have HttpRoute) + DestinationRule 

### Configure Istio auto inject for everything in the namespace

`kubectl label namespace default istio-injection=enabled`


## Bibliography 

[Traffic management in Istio](https://istio.io/latest/docs/concepts/traffic-management/)

[Ingress Gateway into cluster for north/south traffic ](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/)

[Mesh traffic info east/west traffic inside of cluster](https://istio.io/latest/docs/tasks/traffic-management/request-routing/)