### Tools used

1. Go
2. Docker
3. [Skaffold](https://skaffold.dev/)
4. [K3D](https://k3d.io/v5.4.6/)

## 0. Expected tooling to run this project in K3D

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

```127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local grafana.local prometheus.local album-store.local proxy-serivce.local```

### 0.1 K3D Registry info

[K3d Registry info](K3D-registry.md)

## 1. Create K3d Kubernetes Cluster with Internal Registry

```bash
  make k3d-cluster-create;
```

## 2. Build the application in Docker and push Docker image to the  K3D Internal Registry

```bash
  make docker-tag-k3d-registry-album && make docker-tag-k3d-registry-proxy;
```

## 3. Install to cluster: Applications, Observability tooling, Monitoring tooling. 
 
```bash
  make skaffold-dev-k3d;
```

**Note:**

The album-store will not start after printing its version number if OpenTelemetry-collector cannot be reached.

This will mean the liveness probe will fail and the album-store will eventually be in a CrashLoopBackoff state when you get pods.


### 3.1 Debugging Advice  

[Debugging commands for cluster](K3D-Debugging.md)

## 4. Run Some Tests

### 4.1 curl

```bash
curl --insecure --location 'http://album-store.local:8070/albums/'; 
```

### 4.3 Run Test Suite

```bash
make k3d-test;
```

### 4.2 Postman

[Postman files](../test/.)

1. Import the folder `../test`
1. Set Environment to `album-store.local`
1. Open a test in the `Album-Store` collection and run it.

## 5. View the events in the different Services in K3D

[View Jaeger to see spans](http://jaeger.local:8070/search?limit=20&service=album-store)

[View Kubernetes environment](http://k-dashboard.local:8070/)

[Grafana](http://grafana.local:8070/)

[Prometheus](http://prometheus.local:8070/)

## 6. Uninstall from cluster: Applications, Observability tooling, Monitoring tooling.  

Ctr + C on the terminal window where you started `make skaffold-dev`

## 7. Delete Cluster

```bash
  make k3d-cluster-delete;
```

