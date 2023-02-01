# Helm Repo for Project code.

## Helm Charts 

* Album Store
* Proxy Service

## Basic info
Repository with k8s Helm repo.
Basically it's a repository with helm repo.
Check out this example:

- https://github.com/kmzfs/helm-repo-in-github

## Update Helm Repo

```bash
cd helm
#This will create tgz file with chart in charts directory
helm package album-store -d charts 
helm package proxy-service -d charts
helm repo index ./charts #This will create index.yaml file which references album-store.yaml
git add *
git commit -m "album-store helm"
git push
```

## Use Helm Repo
```bash
helm repo add go-gin-opentelemetry 'https://mcarr-and.github.io/go-gin-otelcollector/install/helm/charts'
helm repo update
helm repo list
```
