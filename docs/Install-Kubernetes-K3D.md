## 0. Expected tooling to run this project in K3D

1. Go
2. Docker 
3. Skaffold
4. K3D 
5. local changes to your `/etc/hosts` to use nginx-ingress with your  

```127.0.0.1	localhost k-dashboard.local jaeger.local otel-collector.local```


## 1. Start All Observability & Log Viewing Services
 
```bash
make skaffold-dev;
```

## 2. Start album-store Go/Gin Server with flags set

* `-namespace` kubernetes namespace 
* `-otel-location` can be changed from K3D-Nginx `otel-collector.local`
* `-instance-name` kubernetes instance name (unique name when horizontal scaling)

```bash
make local-start;
```

#### Note: the application will not start without the OpenTelemetry collector running

## 3. Run Some Tests

[Postman Collection](../test/Album-Store.postman_collection.json)

```bash
make local-test;
```

## 4. View the events in the different Services in K3D

[View Jaeger](http://jaeger.local:8070/search?limit=20&service=album-store)

[View K-Dashboard to see Kubernetes environment in a browser](http://k-dashboard:8070/)

TODO Prometheus 

## 5. Stop album-store server & Services  

### 1. Stop Server

`Ctr + C` in the terminal window where go is running. 

### 2. Stop Observability and Log Viewing Services

Ctr + C on the terminal window where you started `make skaffold-dev`