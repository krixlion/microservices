# microservices
*My attempt at creating microservices in Go*

To launch these services on your cluster you just need to clone the repo and apply Kubernetes manifests located in K8s/ directory.
```
git clone https://github.com/Krixlion/microservices

kubectl apply -R -f microservices/k8s/
```
You might need to tweek some settings to match your cloud provider. I've been using GKE.
Service images will be pulled from my public repositories on Docker Hub.
