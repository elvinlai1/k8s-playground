#!/bin/bash

# Install qemu guest agent
sudo apt-get install qemu-guest-agent
sudo systemctl start qemu-guest-agent
sudo systemctl enable qemu-guest-agent

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

sudo reboot now

#Check IP forwarding
sysctl net.ipv4.ip_forward

#Check br_netfilter
lsmod | grep br_netfiltersysctl net.ipv4.ip_forward


# Install Kubeadm kubectl 1.28


sudo apt-get update
# apt-transport-https may be a dummy package; if so, you can skip that package
sudo apt-get install -y apt-transport-https ca-certificates curl gpg


# If the directory `/etc/apt/keyrings` does not exist, it should be created before the curl command, read the note below.
# sudo mkdir -p -m 755 /etc/apt/keyrings

curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg



# This overwrites any existing configuration in /etc/apt/sources.list.d/kubernetes.list
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list



sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

#optional
sudo systemctl enable --now kubelet



# Install Kubevip as HA and Load Balancer
# Declare variable 
export VIP=192.168.30.50
export INTERFACE=ens18

# Get Kube-vip latest version
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




#Configure kubeadm with kubevip as load balancer
sudo kubeadm init --control-plane-endpoint=$VIP:6443 --pod-network-cidr=10.244.0.0/16 --upload-certs


# Follow the outputed instructions to complete the kubeadm init process.

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config


# Install flannel on all control plane
kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml


# Install kube-vip cloud provider
kubectl apply -f https://raw.githubusercontent.com/kube-vip/kube-vip-cloud-provider/main/manifest/kube-vip-cloud-controller.yaml

# Configure kube-vip cloud provider with ip range 
kubectl create configmap --namespace kube-system kubevip --from-literal range-global=192.168.0.200-192.168.0.202


