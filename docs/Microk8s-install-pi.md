
## Dump of notes of what was needed to get a RaspberryPi cluster working. 

### Original instructions
[Ubuntu Microk8s Pi cluster](https://ubuntu.com/tutorials/how-to-kubernetes-cluster-on-raspberry-pi#1-overview)

#### Note
Follow instructions to add a `--worker` node

#### Show nodes in cluster
`microk8s.kubectl get node`

## Post install instructions

### Allow ssh onto Ubuntu machine

This has to be done on each Raspberry Pi as it is not enough from the raspberry Pi Imager to allow ssh alone. 

```bash
 sudo systemctl enable ssh
 ```

### Install Docker & allow running Docker as current user
 ```bash
sudo apt-get install docker.io;
sudo usermod -aG docker $USER
sudo groupadd docker
newgrp docker
```

### allow access to Microk8s as current user
```bash
sudo usermod -aG microk8s $USER
sudo groupadd microk8s
newgrp microk8s
```

### get nodes & enable services for later.

* registry
    * allow images in kubernetes
* ingress
    * allow ingress into cluster
* dns
    * publicly reachable dns???

```
microk8s.kubectl get node
microk8s enable registry ingress dns
 ```


## Opened ports

|purpose                | port number |
|-----------------------|-------------|
| ssh                   | 22 |
| http traffic          | 80 |
| K8s Image registry    | 32000 |
