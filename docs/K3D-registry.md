
## K3d internal registry use 
K3D has a registry you can push locally and then pull into your cluster.


### 0. setup cluster 
This assumes you have use the cluster creation from the root directory of this project. 

```
make k3d-cluster-create
``` 

Which creates a cluster with a registry wired in.

The registry is internally: 

Name `k3d-registry` 

Port `54094`

### 1. build your Docker image 

Supplied Dockerfile which prints `hello world` 

```bash
docker build  . -t forketyfork/hello-world;
```

### 2. tag your image with the repo 

Note: Use `localhost` when building from command line 

```bash
docker tag forketyfork/hello-world:latest localhost:54094/forketyfork/hello-world:v0.1;
```

### 3. push your image with tag into the repository 

Note: Use `localhost` when building from command line

```bash
docker push localhost:54094/forketyfork/hello-world:v0.1;
```

### 4. apply a deployment to use this hello-world image in your k8s cluster 

Note: internally this uses registry url `k3d-registry` 

```bash
 kubectl apply -f helloworld.yaml;
 ```

### 5. get hello world message 

```bash
  kubectl logs $(kubectl get pods -l job-name=hello-world -o jsonpath="{.items[0].metadata.name}")
```

### 6. delete hello-world

Note: internally this uses registry url `k3d-registry`

```bash
 kubectl delete -f helloworld.yaml;
 ```

## Bibliography: 

[K3D Registry creation](https://k3d.io/v4.4.8/usage/commands/k3d_registry_create/)

[Docker build locally instructions](https://medium.com/swlh/how-to-run-locally-built-docker-images-in-kubernetes-b28fbc32cc1d)