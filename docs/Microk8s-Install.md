## 0. Expected tooling to run this project in K3D

1. Go
2. Docker 
3. [Skaffold](https://skaffold.dev/)
4. [Microk8s](https://microk8s.io/)

WIP not currently running. 

TODO improve ths.

## 0. ./kube/config with Microk8s

### 0.1 Setup Kube Config

[Setup kube config](Microk8s-Kube-Config-Setup.md)

### 0.2. Changes to Hosts file

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

```127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local album-store.local```

## 2. Start All Observability & Log Viewing Services
 
```bash
make skaffold-dev-microk8s;
```

## 3. Start album-store Go/Gin Server with flags set

```bash
make local-start-k3d-grpc;
```

#### Note: the application will not start without the OpenTelemetry collector running

## 4. Run Some Tests

### 4.1 Run from command line 

```bash
make local-test;
```

### 4.2 Use Postman

[Postman files](../test/postman_collection.json)

1. Import the collection and environment into your postman
1. Set Environment to `album-store.local`
1. Open a test in the `Album-Store` collection and run it.

## 5. View the events in the different Services in K3D

[View Jaeger](http://jaeger.local:8070/search?limit=20&service=album-store)

[View K-Dashboard to see Kubernetes environment in a browser](http://k-dashboard:8070/)

TODO Prometheus 

## 6. Stop album-store server & Services  

### 1. Stop Server

`Ctr + C` in the terminal window where go is running. 

### 2. Stop Observability and Log Viewing Services

Ctr + C on the terminal window where you started `make skaffold-dev`

## 7. Delete K3D Cluster

```bash
make k3d-cluster-delete
```
