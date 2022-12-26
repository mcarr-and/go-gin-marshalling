## 0. Expected tooling to run this project in K3D

1. Go
2. Docker 
3. Skaffold
4. Microk8s
5. local changes to your `/etc/hosts` to use nginx-ingress with your  

```127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local```

## 1. Add the services 

```bash
microk8s enable jaeger dns prometheus istio dashboard
```
## 2. Add Microk8s to your ~/.kube/config  

### 2.1 get the microk8s kube-config

```bash
microk8s kubectl config view --raw 
```

### 2.2 edit your ~/.kube/config

open your config for kubernetes in your favourite text editor `~/.kube/config`


### 2.3. add your microk8s-cluster to clusters 
to your `clusters:` yaml add 

```yaml
- cluster:
    certificate-authority-data: YOUR_CERT_DATA
    server: https://192.168.205.3:16443
  name: microk8s-cluster
```
### 2.4. add your microk8s user

add to `users:`

```yaml
- name: admin@microk8s
  user:
    token: YOUR_TOKEN
```

### 2.5. add your microk8s context

to `contexts:` add

```yaml
- context:
    cluster: microk8s-cluster
    user: admin@microk8s
  name: microk8s

```

### 2.6. Save your ~/.kube/config

save your changes to `~/.kube/config`  

## 3. Start All Observability & Log Viewing Services
 
```bash
make skaffold-dev-microk8s;
```

## 3. Start album-store Go/Gin Server with flags set

* `-namespace` kubernetes namespace 
* `-instance-name` kubernetes instance name (unique name when horizontal scaling)
* `-otel-location` can be changed from K3D-Nginx `otel-collector.local`

```bash
make local-start-k3d-grpc;
```

#### Note: the application will not start without the OpenTelemetry collector running

## 4. Run Some Tests

[Postman Collection](../test/Album-Store.postman_collection.json)

```bash
make local-test;
```

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
