Install Microk8s on a Mac computer.


WIP - trying to install ArgoCD makes the cluster no longer work even with a 40gb disk and 4gb of ram. 

https://microk8s.io/docs/install-multipass


## 1. Install multipass for OSX

install:
* multipass (for use in installing the microk8s instance correctly)

```bash
brew install multipass microk8s;
```

## 2. fix the microk8s multipass image

```bash
microk8s stop;
multipass stop microk8s-vm;
multipass delete microk8s-vm;
multipass purge;
multipass launch --name microk8s-vm --memory 4G --disk 40G;
```

## 3. Show microk8s image

```bash
multipass list;
```

### 3.1 expected result

There is not a internal IP shown so the microk8s image is not running. 

```
microk8s-vm             Running           192.168.XXX.X    Ubuntu 18.04 LTS
```

## 4. Install microk8s

## 4.1 get a shell on the microk8s image

From your host machine shell onto the microk8s image.

```bash
multipass shell microk8s-vm;

```

## 4.2 show microk8s install information.

```bash
snap info microk8s;
```

### 4.2.1 expected output

There are no `commands:` or `services:` for microk8s image is not installed.

```
name:      microk8s
summary:   Kubernetes for workstations and appliances
publisher: Canonical✓
store-url: https://snapcraft.io/microk8s
contact:   https://github.com/canonical/microk8s
license:   Apache-2.0
description: |
  MicroK8s is a small, fast, secure, certified Kubernetes distribution that installs on just about
  any Linux box. It provides the functionality of core Kubernetes components, in a small footprint,
  scalable from a single node to a high-availability production multi-node cluster. Use it for
  offline developments, prototyping, testing, CI/CD. It's also great for appliances - develop your
  IoT apps for K8s and deploy them to MicroK8s on your boxes.
snap-id: EaXqgt1lyCaxKaQCU349mlodBkDCXRcg
channels:
  1.26/stable:           v1.26.1         2023-02-07 (4595) 176MB classic  ...

...
  ```

## 4.3 install microk8s image

Substitute 1.26/stable for the latest version of stable from the above command.

```bash
sudo snap install microk8s --classic --channel=1.26/stable;
```

### 4.4 confirm configuration

```bash
snap info microk8s;
```

#### 4.4.1 expected output 

`services` and `commands` are now available.

```
name:      microk8s
summary:   Kubernetes for workstations and appliances
publisher: Canonical✓
store-url: https://snapcraft.io/microk8s
contact:   https://github.com/canonical/microk8s
license:   Apache-2.0
description: |
  MicroK8s is a small, fast, secure, certified Kubernetes distribution that installs on just about
  any Linux box. It provides the functionality of core Kubernetes components, in a small footprint,
  scalable from a single node to a high-availability production multi-node cluster. Use it for
  offline developments, prototyping, testing, CI/CD. It's also great for appliances - develop your
  IoT apps for K8s and deploy them to MicroK8s on your boxes.
commands:
  - microk8s.add-node
  - microk8s.addons
  - microk8s.cilium
  - microk8s.config
  - microk8s.ctr
  - microk8s.dashboard-proxy
  - microk8s.dbctl
  - microk8s.disable
  - microk8s.enable
  - microk8s.helm
  - microk8s.helm3
  - microk8s.images
  - microk8s.inspect
  - microk8s.istioctl
  - microk8s.join
  - microk8s.kubectl
  - microk8s.leave
  - microk8s.linkerd
  - microk8s
  - microk8s.refresh-certs
  - microk8s.remove-node
  - microk8s.reset
  - microk8s.start
  - microk8s.status
  - microk8s.stop
  - microk8s.version
services:
  microk8s.daemon-apiserver-kicker: simple, enabled, active
  microk8s.daemon-apiserver-proxy:  simple, enabled, inactive
  microk8s.daemon-cluster-agent:    simple, enabled, active
  microk8s.daemon-containerd:       notify, enabled, active
  microk8s.daemon-etcd:             simple, enabled, inactive
  microk8s.daemon-flanneld:         simple, enabled, inactive
  microk8s.daemon-k8s-dqlite:       simple, enabled, active
  microk8s.daemon-kubelite:         simple, enabled, active
snap-id:      EaXqgt1lyCaxKaQCU349mlodBkDCXRcg
tracking:     1.26/stable
refresh-date: today at 15:49 GMT
channels:
  1.26/stable:           v1.26.1         2023-02-07 (4595) 176MB classic

...
```
## 5. exit from multipass

```bash
exit;
```

## 6. start microk8s 

```bash
microk8s start; 
```

## 7 Enable community addons for microk8s

### 7.1 enable community 
```bash
microk8s enable community;
```

### 7.2 expected results
```bash
microk8s status;
```

### 7.2.1 expected results

```
microk8s is running
high-availability: no
  datastore master nodes: 127.0.0.1:19001
  datastore standby nodes: none
addons:
  enabled:
    community            # (core) The community addons repository
    ha-cluster           # (core) Configure high availability on the current node
    helm                 # (core) Helm - the package manager for Kubernetes
    helm3                # (core) Helm 3 - the package manager for Kubernetes
  disabled:
    argocd               # (community) Argo CD is a declarative continuous deployment for Kubernetes.
    cilium               # (community) SDN, fast with full network policy
    dashboard-ingress    # (community) Ingress definition for Kubernetes dashboard
    fluentd              # (community) Elasticsearch-Fluentd-Kibana logging and monitoring
    gopaddle-lite        # (community) Cheapest, fastest and simplest way to modernize your applications
    inaccel              # (community) Simplifying FPGA management in Kubernetes
    istio                # (community) Core Istio service mesh services
    jaeger               # (community) Kubernetes Jaeger operator with its simple config
    kata                 # (community) Kata Containers is a secure runtime with lightweight VMS
    keda                 # (community) Kubernetes-based Event Driven Autoscaling
    knative              # (community) Knative Serverless and Event Driven Applications
    kwasm                # (community) WebAssembly support for WasmEdge (Docker Wasm) and Spin (Azure AKS WASI)
    linkerd              # (community) Linkerd is a service mesh for Kubernetes and other frameworks
    multus               # (community) Multus CNI enables attaching multiple network interfaces to pods
    nfs                  # (community) NFS Server Provisioner
    ondat                # (community) Ondat is a software-defined, cloud native storage platform for Kubernetes.
    openebs              # (community) OpenEBS is the open-source storage solution for Kubernetes
    openfaas             # (community) OpenFaaS serverless framework
    osm-edge             # (community) osm-edge is a lightweight SMI compatible service mesh for the edge-computing.
    portainer            # (community) Portainer UI for your Kubernetes cluster
    sosivio              # (community) Kubernetes Predictive Troubleshooting, Observability, and Resource Optimization
    traefik              # (community) traefik Ingress controller
    trivy                # (community) Kubernetes-native security scanner
    cert-manager         # (core) Cloud native certificate management
    dashboard            # (core) The Kubernetes dashboard
    dns                  # (core) CoreDNS
    gpu                  # (core) Automatic enablement of Nvidia CUDA
    host-access          # (core) Allow Pods connecting to Host services smoothly
    hostpath-storage     # (core) Storage class; allocates storage from host directory
    ingress              # (core) Ingress controller for external access
    kube-ovn             # (core) An advanced network fabric for Kubernetes
    mayastor             # (core) OpenEBS MayaStor
    metallb              # (core) Loadbalancer for your Kubernetes cluster
    metrics-server       # (core) K8s Metrics Server for API access to service metrics
    minio                # (core) MinIO object storage
    observability        # (core) A lightweight observability stack for logs, traces and metrics
    prometheus           # (core) Prometheus operator for monitoring and logging
    rbac                 # (core) Role-Based Access Control for authorisation
    registry             # (core) Private image registry exposed on localhost:32000
    storage              # (core) Alias to hostpath-storage add-on, deprecated

```

## 8. install addons to make cluster to make it more functional


### 8.1 install services

Install services:
* dashboard (kubernetes dashboard web based portal)
* registry
* observability (Prometheus)
* istio
* jaeger

```bash
microk8s enable dashboard registry observability istio;
microk8s enable jaeger;
```

This will take about 15 minutes to run.


### 10.2 show status

```bash
microk8s status;
```

### 10.2.1 expected result

```
microk8s is running
high-availability: no
  datastore master nodes: 127.0.0.1:19001
  datastore standby nodes: none
addons:
  enabled:
    community            # (core) The community addons repository
    argocd               # (community) Argo CD is a declarative continuous deployment for Kubernetes.
    dashboard            # (core) The Kubernetes dashboard
    dns                  # (core) CoreDNS
    ha-cluster           # (core) Configure high availability on the current node
    helm                 # (core) Helm - the package manager for Kubernetes
    helm3                # (core) Helm 3 - the package manager for Kubernetes
    hostpath-storage     # (core) Storage class; allocates storage from host directory
    istio                # (community) Core Istio service mesh services
    jaeger               # (community) Kubernetes Jaeger operator with its simple config    
    metrics-server       # (core) K8s Metrics Server for API access to service metrics
    observability        # (core) A lightweight observability stack for logs, traces and metrics
    registry             # (core) Private image registry exposed on localhost:32000
    storage              # (core) Alias to hostpath-storage add-on, deprecated
  disabled:
    cilium               # (community) SDN, fast with full network policy
    dashboard-ingress    # (community) Ingress definition for Kubernetes dashboard
    fluentd              # (community) Elasticsearch-Fluentd-Kibana logging and monitoring
    gopaddle-lite        # (community) Cheapest, fastest and simplest way to modernize your applications
    inaccel              # (community) Simplifying FPGA management in Kubernetes
    kata                 # (community) Kata Containers is a secure runtime with lightweight VMS
    keda                 # (community) Kubernetes-based Event Driven Autoscaling
    knative              # (community) Knative Serverless and Event Driven Applications
    kwasm                # (community) WebAssembly support for WasmEdge (Docker Wasm) and Spin (Azure AKS WASI)
    linkerd              # (community) Linkerd is a service mesh for Kubernetes and other frameworks
    multus               # (community) Multus CNI enables attaching multiple network interfaces to pods
    nfs                  # (community) NFS Server Provisioner
    ondat                # (community) Ondat is a software-defined, cloud native storage platform for Kubernetes.
    openebs              # (community) OpenEBS is the open-source storage solution for Kubernetes
    openfaas             # (community) OpenFaaS serverless framework
    osm-edge             # (community) osm-edge is a lightweight SMI compatible service mesh for the edge-computing.
    portainer            # (community) Portainer UI for your Kubernetes cluster
    sosivio              # (community) Kubernetes Predictive Troubleshooting, Observability, and Resource Optimization
    traefik              # (community) traefik Ingress controller
    trivy                # (community) Kubernetes-native security scanner
    cert-manager         # (core) Cloud native certificate management
    gpu                  # (core) Automatic enablement of Nvidia CUDA
    host-access          # (core) Allow Pods connecting to Host services smoothly
    ingress              # (core) Ingress controller for external access
    kube-ovn             # (core) An advanced network fabric for Kubernetes
    mayastor             # (core) OpenEBS MayaStor
    metallb              # (core) Loadbalancer for your Kubernetes cluster
    minio                # (core) MinIO object storage
    prometheus           # (core) Prometheus operator for monitoring and logging
    rbac                 # (core) Role-Based Access Control for authorisation
```