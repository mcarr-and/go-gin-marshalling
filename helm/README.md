# WIP - not working

TODO get this flushed out and running.

## TODO:
* dnsPolicy: ClusterFirst


# Helm Repo
Repository with k8s Helm repo.
Basically it's a repository with helm repo.
Check out this example:
- https://github.com/kmzfs/helm-repo-in-github

## Update Helm Repo
```bash
cd helm
helm package album-store --destination charts #This will create tgz file with chart in charts directory
helm repo index ./charts #This will create index.yaml file which references album-store.yaml
git add *
git commit -m "album-store helm"
git push
```

## Use Helm Repo
```bash
helm repo add album-store 'https://raw.githubusercontent.com/mcarr-and/go-gin-otelcollector/master/helm/charts'
helm repo update
helm repo list
```
