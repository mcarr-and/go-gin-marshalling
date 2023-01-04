### Tools used

1. Go
2. Docker
3. Skaffold
4. K3D

## 0. Expected tooling to run this project in K3D

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

```127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local album-store.local```

### 0.1 K3D Registry info

[K3d Registry info](K3D-registry.md)

## 1. Create K3d Cluster

```bash
make k3d-cluster-create
```

## 2. Build the application and deploy to K3D

```bash
make docker-build;
make k3d-docker-registry;
```

## 3. Start All Observability & Log Viewing Services
 
```bash
make skaffold-dev-k3d;
```

## 4. Deploy Album-Store to Kubernetes 
```bash
make k3d-internal-deploy;
```

Service starts but is not available from external.

This fails to show the JSON payload you get a 503.

```bash
curl -v http://album-store.local:8070/albums/ GET -H "Content-Type: application/json" -H "Host: http://album-store.local:8070
```

[Debugging commands for cluster](Debugging.md)

#### Note: the application will hang after printing its version number if  OpenTelemetry collector is not running

## 5. Run Some Tests

[Postman Collection](../test/Album-Store.postman_collection.json)

```bash
make local-test;
```

## 6. View the events in the different Services in K3D

[View Jaeger](http://jaeger.local:8070/search?limit=20&service=album-store)

[View K-Dashboard to see Kubernetes environment in a browser](http://k-dashboard:8070/)

TODO Prometheus 

## 7. Stop album-store server & Services  

### 1. Stop Server

`Ctr + C` in the terminal window where go is running. 

### 2. Stop Observability and Log Viewing Services

Ctr + C on the terminal window where you started `make skaffold-dev`

## 8. Delete K3D Cluster

```bash
make k3d-cluster-delete
```
