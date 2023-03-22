## 0. Expected tooling to run this project in K3D

1. Go
2. Docker 

## 1. Install microk8s for OSX

`brew install microk8s`

## 2. Install microk8s on the machine image.

### 2.1. get a shell on the microk8s image

`multipass shell microk8s-vm;` 

### 2.2 

`snap info microk8s;`

#### 2.2.1 expected output

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


This means that microk8s in not installed as there are no `commands:` or `services:`

### 2.3 install microk8s

Substitute 1.26/stable for teh latest version of stable from the above command.

`sudo snap install microk8s --classic --channel=1.26/stable;`

### 2.3.1 confirm configuration

`snap info microk8s;`

#### 2.3.1.1 expected output 

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

`services` and `commands` are now available.

### 2.4 make host changes.

#### 2.4.1 exit out of the multipasss shell 

`exit;`

#### 2.4.2 add microk8s config to local expected directory

```bash
mkdir ~/.microk8s;
microk8s config > ~/.microk8s/config;
```

## 2.5 Confirm it is working

```bash
microk8s kubectl get pods;
```

### expected results

`No resources found in default namespace.`