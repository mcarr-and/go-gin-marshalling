### Tools used

1. Go
2. Docker
3. [Skaffold](https://skaffold.dev/)
4. [K3D](https://k3d.io/v5.4.6/)

## 0. Expected tooling to run this project in K3D

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

```127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local grafana.local prometheus.local album-store.local```

### 0.1 K3D Registry info

[K3d Registry info](K3D-registry.md)

## 1. Create K3d Kubernetes Cluster with Internal Registry

```bash
make k3d-cluster-create
```

## 2. Build the application in Docker and push Docker image to the  K3D Internal Registry

```bash
make docker-build-album;
make docker-tag-k3d-registry;
```

## 3. Start All Observability & Log Viewing Services
 
```bash
make skaffold-dev-k3d;
```

## 4. Deploy Album-Store to the K3D Kubernetes Cluster

This will deploy 3 replicas of album-store into the cluster. 

You will see different instance names in the Jaeger Process for the 3 pods.

```bash
make k3d-album-deploy-deployment;
```

**Note: the application will hang after printing its version number if  OpenTelemetry collector is not running**

### 4.1 Debugging Advice  

[Debugging commands for cluster](K3D-Debugging.md)

## 5. View the events in the different Services in K3D

[View Jaeger](http://jaeger.local:8070/search?limit=20&service=album-store)

[View Kubernetes environment](http://k-dashboard:8070/)

## 6. Run Some Tests

### 6.1 curl

```bash
curl --insecure --location 'http://album-store.local:8070/albums/'; 
```

### 6.3 Run Test Suite

```bash
make k3d-test;
```

### 6.2 Postman

[Postman files](../test/.)

1. Import the folder `../test`
1. Set Environment to `album-store.local`
1. Open a test in the `Album-Store` collection and run it.

## 7. Stop album-store server & Services  

Ctr + C on the terminal window where you started `make skaffold-dev`

## 8. Delete Album-Store

```bash
make k3d-album-undeploy-deployment;
```

## 9. Delete Cluster

```bash
make k3d-cluster-delete
```

