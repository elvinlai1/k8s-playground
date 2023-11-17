# k8s-playground
Personal repository revolving around kubernetes (k8s). 
My **on-premise** Kuberrnetes deployment configuration is detailed [here](./MyKubernetesConfiguration.md). 
<br>
<br>
Kubernetes was never designed to be deployed on-premise, but thankfully there are third party components that mimic and replace the need for dedicated cloud components.  

## File structure

```
.
├── go-endpoint                     # Standalone REST API test application
├── go-endpoint-mongodb             # Basic REST API that interacts with mongoDB
│   └── go-endpoint-mongodb.yaml        # Kubernetes deployment config 
│   └── test-request.http               # Predefined requests to test endpoints 
├── mongo                           # Kubernetes MongoDB and Mongo-express deployment configurations
│   └──  docker-compose.yml             # For local device docker testing
└── reference_yaml                  # Place to store yaml configs for reference
├── MyKubernetesConfiguration.md
└── README.md
```

## Table of Contents
- [k8s-playground](#k8s-playground)
    - [Goal(s)](#goals)
    - [Work In Progress](#work-in-progress)
- [What is Kubernetes?](#what-is-kubernetes)
- [My Configured Kubernetes Overview](#the-result-overview)
    - [Visual Topology](#visual-topology)
    - [Nodes information output](#nodes-information-output)
- [Basic Kubernetes feature testing](#basic-kubernetes-feature-testing) 
    - [Kill pods](#kill-pods)
    - [kill nodes](#kill-nodes)
        - [kill worker node](#kill-a-worker-node)
        - [kill control plane node](#kill-a-control-plane)
        - [kill both worker and control plane node](#kill-both-a-control-plane-and-worker-node)
- [Acknowledgements](#acknowledgements)
    - [Variations of k8s](#variations-of-k8s)
- [Appendix](#appendix)

> Acknowledgement section covers any extra knowledge, articles and random thoughts. 

## [Goal(s)](#table-of-contents) 
- Deploy a stacked HA Kubernetes topology cluster via **kubeadm**  
- Create custom manifest to deploy custom `go-endpoint-mongodb` pods to interact with mongoDB 
- Conduct basic stress and load tests to see how kubernete's features work
- Document the knowledge and configurations in a self-hosted environment for future reference

## [Work In Progress](#table-of-contents) 
- Load and Stress testing to understand kubernete's automated system
- Prometheus configuration monitoring 
- Configuring Storage Orchestration
    - Persistent Volume Claim (PVC) --> Persistent Volume (PV)
- Proper Nginx Ingress controller and TLS/SSL termination  
- Automated CI/CD setup with Jenkins 

# [What is Kubernetes?](#table-of-contents)

<br>
<br>

# [My Configured Kubernetes Overview](#table-of-contents) 
To show case what my current configuration is capable of 

## Visual Topology 
"Image"

## Nodes information output


# [Basic Kubernetes feature testing](#table-of-contents)

## Self-healing
(gif)
### kill pods
(gif)

### Kill nodes

#### Kill a worker node 
(gif)

#### Kill a control plane 
(gif)

#### Kill both worker and control plane node
(gif)




# [Acknowledgements](#table-of-contents)

As **kubeadm** is just a tool to create a minimum viable cluster. More in-depth configuration can be done as shown in [kelseyhightower ](https://github.com/kelseyhightower/kubernetes-the-hard-way ) and [mmumshad](https://github.com/mmumshad/kubernetes-the-hard-way) repository. 


## Variations of k8s 
- minikube 
- kind
- k3s 
- microk8s

## k8s components
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



[TalosOS](https://www.talos.dev/)



# [Appendix](#table-of-contents)

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

