# On-Premise Kubernetes Cluster Config

> Personal local cluster for testing and learning
- ARP Mode
- High Availability and Load balanced for control plane and services 
- Virtual IP for leader-election
- On-Premise Cloud Provider component substitution

---

<details>
<summary> Network Configuration </summary>

- VLAN: `30`
- Subnet: `192.168.30.0/24`
- DNS/Gateway: `192.168.30.1`

</details>

---

<details>
<summary> Node / VM Configuration </summary>

Proxmox VE:
- Control Plane node:
	- Ubuntu 22.04 
	- 4 cores
	- 4GB ram 
	- 40GB storage
	- VLAN `30`
- Worker node: 
	- Ubuntu 22.04 
	- 2 cores
	- 4GB ram 
	- 40GB storage
	- VLAN `30`
<details>
<summary> Qemu Guest Agent (optional) </summary>

```bash
sudo apt-get install qemu-guest-agent
sudo systemctl start qemu-guest-agent
sudo systemctl enable qemu-guest-agent
```
</details>
</details>

---

<details open>
<summary>Package Versions</summary>

- kubectl kubeadm kubelet `v1.28` (specified version due to [#issue](https://github.com/kube-vip/kube-vip/issues/684))
- kube-vip `v0.8.2` (or latest release)
- kube-vip cloud provider `v0.0.10` (or latest release)
- flannel [CNI] `v0.25.6` (or latest release)
- containerd [CRI] `v1.5.0-2.1` (or latest release)

</details>

---

<details>
<summary> Node Configuration Scripts [WiP]</summary>

Control Plane Node Script 
```
curl 
```

Worker Node Script 
```
curl
```

</details>

---

## Config Steps

### 1. Required configurations prior to installing Kubernetes 

<details open>

```bash
# Disable Firewall
sudo ufw disable

# Install CRI and enable cgroup driver
sudo apt install containerd -y
sudo mkdir /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml
sudo sed -i 's/SystemdCgroup = .*/SystemdCgroup = true/' /etc/containerd/config.toml

# Disable Swap
sudo sed -i 's/^\/swap\.img/#\/swap.img/' /etc/fstab

# Enable bridging
sudo sed -i 's/#net\.ipv4\.ip_forward=1/net.ipv4.ip_forward=1/' /etc/sysctl.conf

# Enable br_netfilter 
echo "br_netfilter" | sudo tee /etc/modules-load.d/k8s.conf > /dev/null
```
</details>


### 2. Reboot
`sudo reboot now`

### 3. Confirm configs

<details open>

```bash
sudo cat /etc/containerd/config.toml | grep "SystemdCgroup"
sudo cat /etc/fstab | grep "swap"
sudo cat /etc/sysctl.conf | grep "net.ipv4.ip_forward"
sudo cat /etc/modules-load.d/k8s.conf
```
</details>


### 4. Installing kubeadm, kubelet and kubectl

<details open>

```bash
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gpg

# Get 1.28 package keys
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

sudo systemctl enable --now kubelet
```
</details>


### 5. Install Kube-vip as Static Pods

<details open>

```bash
# Install Kubevip as HA and Load Balancer
export VIP=192.168.30.50 # Specified Floating Virtual IP out of DHCP range
export INTERFACE=ens18 # VM network interface name
# Get Kube-vip latest release
KVVERSION=$(curl -sL https://api.github.com/repos/kube-vip/kube-vip/releases | jq -r ".[0].name")

# Set non-persistent alias injected with "sudo"
alias kube-vip="sudo ctr image pull ghcr.io/kube-vip/kube-vip:$KVVERSION; sudo ctr run --rm --net-host ghcr.io/kube-vip/kube-vip:$KVVERSION vip /kube-vip"


kube-vip manifest pod \
    --interface $INTERFACE \
    --address $VIP \
    --controlplane \
    --services \
    --arp \
    --leaderElection | sudo tee /etc/kubernetes/manifests/kube-vip.yaml
```

</details>


### 6. Initialize First Control Plane Node

<details open>

```bash
sudo kubeadm init --control-plane-endpoint=$VIP:6443 --pod-network-cidr=10.244.0.0/16 --upload-certs
```

#### Flags explained 
`--control-plane-endpoint` Virtual IP that handles all kubernetes API communication

`--pod-network-cidr` specified as flannel uses 10.244.0.0/16 as default podCIDR range

`--upload certs` Upload control-plane certificates to the kubeadm-certs Secret necessary to join nodes

</details>


### 7. Post `kubeadm init`

<details open>


Test load balancer status
`nc -v 192.168.30.50 6443`

<details>
<summary> Watch for printed out instructions below</summary>

```bash
# Configure permissions to cluster
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

```bash
# Example Control Plane Join command
kubeadm join 192.168.30.50:6443 --token xb4ppx.axq77dn1fg9sy9ow --discovery-token-ca-cert-hash sha256:cdf1655ac588eea7eedf18709747a0980d9250456bf3155f19cb848c5c6cbc59 --control-plane --certificate-key <generated cert hash> 
```

```bash
# Create join tokens and cert hash
kubeadm init phase upload-certs --upload-certs 
kubeadm token create --print-join-command
```


</details>

<br>

> If node is still in `notReady` state for a few minutes
`sudo systemctl restart containerd` or `ERROR: cni plugin not intialized`

</details>


### 8. Install the Container Network Interface (CNI)

<details open>

[Flannel](https://github.com/coreos/flannel)  L3 Fabric for Kubernetes

```bash
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
```
</details>


### 9. Install On-premise Kube-vip Cloud Provider 

<details open>

> Provides external IP for LoadBalancer services from a specfied CIDR-range

[Kube-vip Cloud Provider Documentation](https://kube-vip.io/docs/usage/cloud-provider/)

```bash
# Install cloud provider
kubectl apply -f https://raw.githubusercontent.com/kube-vip/kube-vip-cloud-provider/main/manifest/kube-vip-cloud-controller.yaml
```

```bash
# Create configmap of IP ranges for cloud provider
kubectl create configmap --namespace kube-system kubevip --from-literal range-global=<IP Start>-<IP End>
```
</details>

### 10. Deploy Pod with Service
<details open>


#### Personal docker image deployments
> go-endpoint-deployment
```
kubectl apply -f https://github.com/elvinlai1/k8s-playground/blob/main/go-endpoint/go-endpoint.yaml
```

>go-endpoint-mongodb-deployment
```
kubectl apply -f https://github.com/elvinlai1/k8s-playground/blob/main/go-endpoint/go-endpoint-mongodb-deployment.yaml
```

</details>
