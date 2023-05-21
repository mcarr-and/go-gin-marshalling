
## Dump of notes of what was needed to get a RaspberryPi cluster working. 

WIP - DNS not resolving urls 

TODO: 
* DNS resolution is still a work in progress.

## /etc/hosts on your local machine

local changes to your `/etc/hosts` to use nginx-ingress with the k3d cluster.

change `192.168.XX.XX` to the IP Address of the raspberrypi that is the control plane.

```192.168.XX.XX	k-dashboard.local jaeger.local otel-collector.local grafana.local prometheus.local kiali.local album-store.local proxy-serivce.local registry.local```


## Original instructions
[Ubuntu Microk8s Pi cluster](https://ubuntu.com/tutorials/how-to-kubernetes-cluster-on-raspberry-pi#1-overview)

#### Note
Follow instructions to add a `--worker` node

#### Show nodes in cluster
`microk8s.kubectl get node`


# Post install instructions

## Turn on services (Control Plane and all Worker Nodes)

```bash
  sudo systemctl enable ssh;
  # sudo ufw enable; # disabled until all ports are known 
 ```

## Install Docker & allow running Docker as current user (Control Plane and all Worker Nodes)
 ```bash
  sudo apt-get install docker.io;
  sudo usermod -aG docker $USER;
  sudo groupadd docker;
  newgrp docker;
```

## allow access to Microk8s as current user (Control Plane and all Worker Nodes)
```bash
  sudo usermod -aG microk8s $USER;
  sudo groupadd microk8s;
  newgrp microk8s;
```

## Create Docker Registry (Control Plane) 

```bash
sudo bash -c "cat > /etc/docker/daemon.json << EOF 
{
  \"insecure-registries\" : [\"localhost:32000\"]
}
EOF";
sudo systemctl restart docker;
```

## Enable microk8s services that are needed(Control Plane)

```bash
microk8s.kubectl get nodes;
microk8s enable community;
microk8s enable hostpath-storage; # needed for Jaeger to mount a volume
microk8s enable dns; # needed for default DNS resolution
microk8s enable registry # allow saving of local docker images 
 ```

## Apple - Allow Rancher Desktop - Docker to use insecure registry

The following needs to be run each time the computer is restarted or Rancher Desktop is updated.

This allows docker to push to the microk8s registry. 

```bash
LIMA_HOME="$HOME/Library/Application Support/rancher-desktop/lima" "/Applications/Rancher Desktop.app/Contents/Resources/resources/darwin/lima/bin/limactl" shell 0;
sudo sed -i 's/DOCKER_OPTS=\"\"/DOCKER_OPTS=\"--insecure-registry=registry.local:32000\"/g' /etc/conf.d/docker;
sudo service docker restart;
```

### Error Message if not set
You will get an error message that states that an http registry is not allowed:

```
The push refers to repository [registry.local:32000/album-store]
Get https://registry.local:32000/v2/: http: server gave HTTP response to HTTPS client
```

## DNS Fix /etc/resolv.conf (Control Plane and all Worker Nodes)

```bash
sudo sed -i 's/#Domains=/Domains=svc.cluster.local cluster.local/g' /etc/systemd/resolved.conf;
sudo systemctl restart systemd-resolved;
```

### DNS issues

This allows the searching of internal domains `svc.cluster.local` and `cluster.local` in the K8s cluster.

The below change will add those 2 entries in your `/etc/resolv.conf` and it will end with the name of your router.

It will try to DNS resolve inside the cluster first before going external.

### Istio DNS issues
We need this as istio-gateway needs istiod. Gateway should search the cluster to find istiod.

You need to have istiod running before installing istio-gateway.

In skaffold you have to add `wait: true` in the istiod before installing istio-gateway.

## istio-gateway fixes (Control Plane and all Worker Nodes)

One needs to add entries to give a clusterIP to the isto 

```bash
export CLUSTER_DNS=$(kubectl -nkube-system get svc/kube-dns -o jsonpath="{.spec.clusterIP}");
sudo bash -c "cat >> /var/snap/microk8s/current/args/kubelet << EOF
--cluster-domain=cluster.local
--cluster-dns=$CLUSTER_DNS
EOF";
sudo snap stop microk8s; sudo snap start microk8s;
```

## Install tooling and applications

Replace XXX.XXX.XX.XX; with your ip address of one of the worker nodes.

```bash
  export ISTIO_GATEWAY_EXTERNAL_IP=XXX.XXX.XX.XX;
  make skaffold-dev;
```


## Firewall

This is still WIP so not used. 

Istio-Gateway was not finding istiod and failing to launch if ufw rules are set. 

### Firewall Rules (Control Plane)

```bash
#sudo ufw enable;
#sudo ufw allow 22/tcp; # ssh
#sudo ufw allow 53/tcp; # dns
#sudo ufw allow 80/tcp; # http
#sudo ufw allow 443/tcp; # k8s metrics port
#sudo ufw allow 2379:2380/tcp; # etcd  
#sudo ufw allow 6443/tcp; # kubectl api port
#sudo ufw allow 8001/tcp; # istio liveness
#sudo ufw allow 8080/tcp; # istio debug port
#sudo ufw allow 8081/tcp; # kube telemetry port
#sudo ufw allow 9153/tcp; # prometheus metrics port
#sudo ufw allow 9100/tcp; # prometheus node exporter 
#sudo ufw allow 9091/tcp; # prometheus push gateway 
#sudo ufw allow 9093/tcp; # prometheus alert manager  
#sudo ufw allow 10250:10252/tcp; # kubelet, kube-schedule, kube-controller 
#sudo ufw allow 10255/tcp; # kubelet
#sudo ufw allow 10257/tcp; # kube-scheduler-manager
#sudo ufw allow 10259/tcp; # kube-scheduler
#sudo ufw allow 16443/tcp; # kubectl
#sudo ufw allow 25000/tcp; # worker nodes to join the node group
#sudo ufw allow 15000/tcp; # istio envoy admin
#sudo ufw allow 15001/tcp; # istio envoy outbound
#sudo ufw allow 15004/tcp; # istio envoy debug
#sudo ufw allow 15006/tcp; # istio envoy inbound
#sudo ufw allow 15008:15009/tcp; #istio hbone traffic
#sudo ufw allow 15010/tcp; # istio xds and ca
#sudo ufw allow 15012/tcp; # istio xds and ca
#sudo ufw allow 15014/tcp; # istiod control plane monitoring
#sudo ufw allow 15017/tcp; # istiod webhook container port
#sudo ufw allow 15020/tcp; # istio status
#sudo ufw allow 15021/tcp; # istio health check
#sudo ufw allow 15053/tcp; # istio dns capture
#sudo ufw allow 15090/tcp; # istio envoy prometheus
#sudo ufw allow 32000/tcp; # docker registry
```

### Firewall Rules (Worker Node)

```bash
#sudo ufw enable;
#sudo ufw allow 22/tcp; # ssh
#sudo ufw allow 53/tcp; # dns
#sudo ufw allow 80/tcp; # http
#sudo ufw allow 443/tcp; # k8s metrics port
#sudo ufw allow 2379:2380/tcp; # etcd 
#sudo ufw allow 6443/tcp; # kubectl api port 
#sudo ufw allow 8001/tcp; # istio liveness?
#sudo ufw allow 8080/tcp; # istio debug port
#sudo ufw allow 9153/tcp; # prometheus metrics port 
#sudo ufw allow 9100/tcp; # prometheus node exporter 
#sudo ufw allow 9091/tcp; # prometheus push gateway 
#sudo ufw allow 9093/tcp; # prometheus alert manager 
#sudo ufw allow 10250:10252/tcp; # kubelet, kube-schedule, kube-controller
#sudo ufw allow 10255/tcp;
#sudo ufw allow 10259/tcp; # kube-scheduler
#sudo ufw allow 10257/tcp; # kube-scheduler-manager
#sudo ufw allow 15000/tcp; #istio envoy admin
#sudo ufw allow 15001/tcp; #istio envoy outbound
#sudo ufw allow 15004/tcp; #istio envoy debug
#sudo ufw allow 15006/tcp; #istio envoy inbound
#sudo ufw allow 15008:15009/tcp; #istio hbone traffic
#sudo ufw allow 15010/tcp; # istio xds and ca
#sudo ufw allow 15012/tcp; # istio xds and ca
#sudo ufw allow 15014/tcp; # istiod control plane monitoring
#sudo ufw allow 15017/tcp; # istiod webhook container port
#sudo ufw allow 15020/tcp; # istio status
#sudo ufw allow 15021/tcp; # istio health check
#sudo ufw allow 15053/tcp; # istio dns capture
#sudo ufw allow 15090/tcp; # istio envoy prometheus
#sudo ufw allow 16443/tcp; # kubectl 
#sudo ufw allow 25000/tcp; # worker nodes to join the node group
#sudo ufw allow 30000:32767/tcp; # allow external traffic
```


## Bibliography

[resolv.conf issues and cluster resolution for istio base](https://rtfm.co.ua/en/kubernetes-load-testing-and-high-load-tuning-problems-and-solutions/#php_network_getaddresses_getaddrinfo_failed_%D0%B8_DNS)

[Microk8s registry](https://microk8s.io/docs/registry-built-in)

[Rancher Desktop fix docker for insecure registries](https://github.com/rancher-sandbox/rancher-desktop/discussions/1477#discussioncomment-2106389)