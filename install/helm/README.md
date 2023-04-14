# Helm Repo for Project code.

## Helm Charts 

* Album Store
* Proxy Service
* Grafana Dashboards via Configmaps. Inspred by [Glenn De Haan's Grafana Helm](https://github.com/glenndehaan/charts/tree/master/charts/grafana-dashboards)

## Basic info
Repository with k8s Helm repo.
Basically it's a repository with helm repo.
Check out this example:

- https://github.com/kmzfs/helm-repo-in-github

## Update Helm Repo

```bash
cd helm
#This will create tgz file with chart in charts directory
helm package album-store -d _charts; 
helm package proxy-service -d _charts;
helm package grafana-dashboards -d _charts;
#This will create index.yaml file which references album-store.yaml and proxy-service.yaml
helm repo index ./_charts; 
git add *;
git commit -m "helm chart updates"
git push
```

## Use Helm Repo
```bash
helm repo add go-gin-opentelemetry 'https://mcarr-and.github.io/go-gin-otelcollector/install/helm/_charts'
helm repo update
helm repo list
```
