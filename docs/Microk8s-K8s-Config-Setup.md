## 0. Expected tooling to run this project in K3D

TODO - config for jager and microk8s & make this work.

WIP not currently running

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
