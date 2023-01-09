### Tools used

1. Go
2. Docker
3. [Skaffold](https://skaffold.dev/)
4. [K3D](https://k3d.io/v5.4.6/)

## 0. Expected tooling to run this project in K3D

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

`proxy-service.local` is included in this configuration if you are following the album-store config changes.

`127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local album-store.local proxy-service.local`

### 0.1 K3D Registry info

[K3d Registry info](../docs/K3D-registry.md)

## 1. Create K3d Kubernetes Cluster with Internal Registry

```bash
  make k3d-cluster-create;
```

## 2. Build the applications in Docker and push Docker image to the  K3D Internal Registry

```bash
  make docker-build-proxy && make docker-tag-k3d-registry;
  make docker-build-album && make docker-tag-k3d-registry-album;
```

## 3. Start All Observability & Log Viewing Services
 
```bash
  make skaffold-dev-k3d;
```

## 4. Deploy Album-Store to the K3D Kubernetes Cluster

This will deploy 3 replicas of album-store into the cluster.

You will see different instance names in the Jaeger Process for the 3 pods.

```bash
  make k3d-album-deploy-deployment && make k3d-proxy-deploy-deployment;
```

**Note: the application will hang after printing its version number if  OpenTelemetry collector is not running**

### 4.1 Debugging Advice  

[Debugging commands for cluster](../docs/K3D-Debugging.md)

## 5. View the events in the different Services in K3D

These will have 2 spans.

* 1 for Proxy-Service
* 1 for Album-Store

[View Jaeger](http://jaeger.local:8070/search?limit=20&service=proxy-service)

[View Kubernetes environment](http://k-dashboard:8070/)

## 6. Run Some Tests

### 6.1 Command line

```bash
curl --insecure --location 'http://proxy-service.local:8070/albums/'; 
```

```bash
make k3d-test;
```

### 6.2 Postman

[Postman files](../test/postman_collection.json)

1. Import the folder `../test`
1. Set Environment to `proxy-service.local`
1. Open a test in the `Album-Store` collection and run it.

### 6.3 Browser

[view proxy-service albums](http://album-service:8070/albums)

## 7. Stop the Services  

Ctr + C on the terminal window where you started `make skaffold-dev`

## 8. Delete Proxy-Service from Kubernetes

```bash
  make k3d-proxy-undeploy-deployment;
```

## 9. Delete ProxyAlbum-Service from Kubernetes

```bash
  make k3d-album-undeploy-deployment;
```