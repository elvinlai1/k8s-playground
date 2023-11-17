# Kubernetes Configuration
This section will provide details regarding hardware and software configurations done to enable an ***on-premise*** **kubeadm** HA stacked etcd topology. 

Proxmox Hypervisor server with each node as a Virtual Machine (VM) 

(kubeadm HA topology stacked etcd images)

Refer to [documentation]( https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/) for more detail on the requirements


## Table of Contents
- [The Setup](#the-setup)
	- [On-Premise Configuraiton]()
	- [Network Configuration](#network-configuration)
	- [Proxmox](#proxmox-hypervisor)
		- [Node VM resource allocation configuration](#node-vm-resource-allocation-configuration)
	- [Configuring HA Kubernetes Cluster](#configuring-ha-kubernetes-cluster)
		- [1. Required configurations prior to installing Kubernetes](#1-required-configurations-prior-to-installing-kubernetes)
		- [2. Reboot](#2-reboot)
		- [3.Confirm Configs](#3-confirm-configs)
		- [4. Installing Kubeadm kubelet and kubectl](#4-installing-kubeadm-kubelet-and-kubectl)
		- [5. Create load balancer for kube apiserver](#5-create-load-balancer-for-kube-apiserver)
			- [5. a. Pre-manifest generation configuration](#5a-pre-manifest-generation-configuration)
		- [6. Initialise master control plane node](#6-initialise-master-control-plane-node)
		- [7.Install the container network interface CNI](#7-install-the-container-network-interface-cni)
		- [8. Deploy pods ](#8-deploy-pods)
		- [9. Expose pods ](#9-expose-pods)

## [On-Premise Configuration](#table-of-contents)

## [Network Configuration](#table-of-contents) 
VLAN: `30`\
Subnet: `192.168.30.0/24`   
DNS/Gateway handled by my router: `192.168.30.1` 

## [Proxmox Hypervisor](#table-of-contents)

### Node VM resource allocation configuration  
- Control node: 
	- Ubuntu 22.04 
	- 4 cores
	- 4GB ram 
	- 20GB storage
	- VLAN `30`
- Worker node: 
	- Ubuntu 22.04 
	- 2 cores
	- 4GB ram 
	- 20GB storage
	- VLAN `30`

## [Configuring HA Kubernetes cluster](#table-of-contents) 
Most of the required configurations are stated by the official kubernetes **kubeadm** installation documentation to reduce any errors during `kubeadm init` startup. 
I have condensed the information corroborating with other sources to simplify the process below. 
<br>
As usual refer to [Documentation](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/) for more detail. 

***Each node must have these base configuration prior to joining the first control-plane.***

### 1. Required configurations prior to installing Kubernetes 
```
# Install Container Runtime Interface and enable cgroup driver
sudo apt install containerd -y
sudo mkdir /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
sudo sed -i 's/SystemdCgroup = .*/SystemdCgroup = true/' /etc/containerd/config.toml

# Disable swap by commenting it out
sudo sed -i 's/^\/swap\.img/#\/swap.img/' /etc/fstab

# Enable Bridging by uncommenting
sudo sed -i 's/#net\.ipv4\.ip_forward=1/net.ipv4.ip_forward=1/' /etc/sysctl.conf

# Enable br_netfilter by appending "br_netfilter to the file"
echo "br_netfilter" | sudo tee /etc/modules-load.d/k8s.conf > /dev/null
```
### 2. Reboot
`sudo reboot now`
### 3. Confirm configs 
```
sudo cat /etc/containerd/config.toml | grep "SystemdCgroup"
sudo cat /etc/fstab | grep "swap"
sudo cat /etc/sysctl.conf | grep "net.ipv4.ip_forward"
sudo cat /etc/modules-load.d/k8s.conf
```
### 4. Installing kubeadm, kubelet and kubectl  
Kubernetes official [documentation](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
) will go through the steps to install the required packages to setup Kubernetes depending on Linux distribution. However, it is expected to change for newer package versions so I decided not to document the commands state-fully.
### 5. Create load balancer for kube-apiserver

This is the crucial step in order to achieve a HA stacked **etcd** k8s cluster.

**Explain stacked etcd**


There are ways to handle this such as a **HAproxy** and **keepalived** combination to provide the required `LoadBalancer` and floating virtual IP (VIP) either on each of the control nodes themselves or a dedicated server. 

**Explain purpose of VIP here**

**Explain HAproxy**

**Explain Metallb** 

**Explain GKE ingress** 


As I am configuring a HA cluster in a self-hosted environment I went with the Layer 2 **Kube-vip** ARP configuration to simplify the process and once configured the master node will auto propagate the configuration to other control-plane nodes as they come online. 

However, I do lose some declared functionality as I am deploying it as a pod
> it can't make use of a variety of Kubernetes functionality (such as the Kubernetes token or ConfigMap resources)

#### 5.a Pre-manifest generation configuration 
Refer to [Documentation](https://kube-vip.io/docs/installation/static/) on configuration and sequential startup details

```
# Declare variable 
export VIP=192.168.30.50
export INTERFACE=ens18

# Get Kube-vip latest version
`KVVERSION=$(curl -sL https://api.github.com/repos/kube-vip/kube-vip/releases | jq -r ".[0].name")`

# Set non-persistent alias injected with "sudo"
alias kube-vip="sudo ctr image pull ghcr.io/kube-vip/kube-vip:$KVVERSION; sudo ctr run --rm --net-host ghcr.io/kube-vip/kube-vip:$KVVERSION vip /kube-vip"
```

```
# Execute alias injected with "sudo"
kube-vip manifest pod \
    --interface $INTERFACE \
    --address $VIP \
    --controlplane \
    --services \
    --arp \
    --leaderElection | sudo tee /etc/kubernetes/manifests/kube-vip.yaml
```

As **Kube-vip** pod is initialised during `kubeadm init` the port designated in your `kubeadm init` configuration will be passed to **Kube-vip** manifest

### 6. Initialise master control plane node
Refer to [Documentation](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/) for more detail regarding `kubeadm init`.
Refer to `kubeadm reset` commands [here](https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-reset/) to reset configuration.

```
sudo kubeadm init --control-plane-endpoint=192.168.30.50:6443 --pod-network-cidr=10.244.0.0/16 --upload-certs
```
#### Current flags 
`--control-plane-endpoint`

`--pod-network-cidr`

`--upload certs` will generate the necessary flags and parameters for `kubeadm join...` to designate and connect the node as either control-plane or worker


#### 6.a Post kubeadm initialisation
Test load balancer status
`nc -v 192.168.30.50 6443` 

Once the command is successful:
- Find the commands printed out to give the current account access to the cluster
- Identify the `kubeadm join ... --control-plane ...` for control plane join configuration
- Identify the `kubeadm join...` for worker node join configuration 

Join tokens expire so create new ones when necessary
`kubeadm token create --print-join-command


### 7. Install the Container Network Interface (CNI)
**Explain the purpose of CNI and issues without it** 

[Flannel](https://github.com/coreos/flannel) is a well-established overlay network and provides a pure Layer 3 network fabric for Kubernetes clusters.

`kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml`


### 8. Deploy pods 




### 9. Expose pods 




