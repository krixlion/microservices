# Disclaimer
**This repo was made as an exercise and should _not_ be used in _production_ environment.**

# microservices
*My attempt at creating microservices entirely written in Go*
<br>
I'm writing my own Eventstore using MongoDB as primary database and I'm using RabbitMQ to handle Pub/Sub communication inbetween services.

# Installation
To launch these services on your cluster you just need to clone the repo and apply Kubernetes manifests located in K8s/ directory.

```
git clone https://github.com/Krixlion/microservices

kubectl apply -R -f microservices/k8s/
```

You might need to tweek some settings to match your cloud provider. I've been using GKE. <br>
Docker images are available in my public repositories on Docker Hub. <br>
You might as well build them yourself from Dockerfiles included in each of services directories.
