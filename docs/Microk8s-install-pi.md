
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
  sudo ufw enable;
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

## Fix /etc/resolv.conf (Control Plane and all Worker Nodes)

```bash
sudo sed -i 's/#Domains=/Domains=svc.cluster.local cluster.local/g' /etc/systemd/resolved.conf;
sudo systemctl restart systemd-resolved;
```

### DNS issues

This allows the searching of internal domains `svc.cluster.local` and `cluster.local` in the K8s cluster.

The below change will add those 2 entries in your `/etc/resolv.conf` and it will end with the name of your router.

It will try to DNS resolve inside the cluster first before going external.

### Istio issues
We need this as istio-gateway needs istiod. Gateway should search the cluster to find istiod.

You need to have istiod running before installing istio-gateway.

In skaffold you have to add `wait: true` in the istiod before installing istio-gateway.

### Firewall Rules (Control Plane)

```bash
sudo ufw enable;
sudo ufw allow 22/tcp; # ssh
sudo ufw allow 80/tcp; # http
sudo ufw allow 6443/tcp; # kubectl
sudo ufw allow 16443/tcp; # kubectl
sudo ufw allow 10250:10252/tcp; # kubelet, kube-schedule, kube-controller 
sudo ufw allow 25000/tcp; # worker nodes to join the node group
sudo ufw allow 10255/tcp; # kubelet
sudo ufw allow 2379:2380/tcp; 
sudo ufw allow 15090/tcp; # prometheus
sudo ufw allow 15012/tcp; #istio health check
sudo ufw allow 15021/tcp; # istio traffic
sudo ufw allow 32000/tcp; # docker registry
sudo ufw allow 53/tcp; # certificate
```

### Firewall Rules (Worker Node)

```bash
sudo ufw enable;
sudo ufw allow 22/tcp; # ssh
sudo ufw allow 10250/tcp;
sudo ufw allow 10255/tcp;
sudo ufw allow 30000:32767/tcp; # allow external traffic
sudo ufw allow 15012/tcp; #istio health check
sudo ufw allow 15021/tcp; # istio traffic
sudo ufw allow 15090/tcp; # prometheus
## unsure
sudo ufw allow 53/tcp; # certificate
sudo ufw allow 80/tcp; # http
```