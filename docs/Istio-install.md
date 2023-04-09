## Prerequsites 

kubectl get crd gateways.gateway.networking.k8s.io &> /dev/null || \
  { kubectl kustomize "github.com/kubernetes-sigs/gateway-api/config/crd?ref=v0.6.1" | kubectl apply -f -; }

## Configure Traffic to pierce the cluster from outside the mesh

### Old way 

Ingress

### Istio Way 

Gateway + HttpRoute ? 

??? When to use VirtualService vs HttpRoute for destination of gateway???

VirtualService is for east west traffic.

### Notes 

Must wait for the Gateway to be ready before creating the HttpRoute

`kubectl wait --for=condition=ready gtw httpbin-gateway`

## Traffic inside cluster = inside mesh

Mesh service = VirtualService(can have HttpRoute) + DestinationRule 

### Configure Istio auto inject for everything in the namespace

`kubectl label namespace default istio-injection=enabled`


## Add external services to Mesh 

Mesh external service = ServiceEntry


## Bibliography 

[Traffic management in Istio](https://istio.io/latest/docs/concepts/traffic-management/)

[Ingress Gateway into cluster for north/south traffic ](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/)

[Mesh traffic info east/west traffic inside of cluster](https://istio.io/latest/docs/tasks/traffic-management/request-routing/)