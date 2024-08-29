# k8s-playground
> Personal repository revolving around kubernetes (k8s).

**On-Premise Configuration**  [here](OnPremK8sCluster.md).

**Cloud Configuration** [WiP]


### Goal(s)
- Deploy a stacked High Availability Kubernetes cluster via **kubeadm**  
- Create custom manifest to deploy custom `go-endpoint-mongodb` pods to interact with mongoDB 
- Conduct basic stress and load tests to see how kubernete's features work
- Document the knowledge and configurations in a self-hosted environment for future reference

### Work In Progress
- Load and Stress testing to understand kubernete's automated system
- Prometheus configuration monitoring
- Configuring Storage Orchestration
    - Persistent Volume Claim (PVC) --> Persistent Volume (PV)
- Proper Nginx Ingress controller and TLS/SSL termination  
- Automated CI/CD setup with Github actions/Jenkins


## References
### Dedicated OS
- [TalosOS](https://www.talos.dev/)

### Variations of k8s 
- minikube 
- kind
- k3s 
- microk8s

### k8s components
- persistent storage 
- namespaces
- rbac 
- containerd
- Flannel (CNI)
- CoreDNS
- Traefik (Ingress)
- Klipper-lb (Service LB)
- Embedded network policy controller
- Embedded local-path-provisioner
- Host utilities (iptables, socat, etc)

### Articles

Explaining the purpose of kubernetes
https://dev.to/thenjdevopsguy/what-problem-is-kubernetes-actually-trying-to-solve-3g1n

Setting up kubernetes in proxmox
https://www.learnlinux.tv/how-to-build-an-awesome-kubernetes-cluster-using-proxmox-virtual-environment/

Article regarding high availability kubernetes using kubeadm
https://medium.com/velotio-perspectives/demystifying-high-availability-in-kubernetes-using-kubeadm-3d83ed8c458b

How Kube-vip replaces both haproxy and keepalived to a degree
https://inductor.medium.com/say-good-bye-to-haproxy-and-keepalived-with-kube-vip-on-your-ha-k8s-control-plane-bb7237eca9fc

Explanation of kube-vip
https://thebsdbox.co.uk/2020/01/02/Designing-Building-HA-bare-metal-Kubernetes-cluster/#Networking-load-balancing

Explanation of multi-cluster kubernetes 
https://traefik.io/glossary/understanding-multi-cluster-kubernetes/

Examples of multi-cluster kubernetes configuration 
https://www.kubecost.com/kubernetes-multi-cloud/kubernetes-multi-cluster/

Kubernetes for Data Engineering 
https://medium.com/@vladimir.prus/kubernetes-for-data-engineering-feba247f7585



