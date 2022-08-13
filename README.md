# Disclaimer
**This repo was made as an exercise and should _not_ be used in _production_ environment.**

## Microservices
*My attempt at creating microservices entirely written in Go*

I'm writing my own Eventstore using MongoDB as primary database and I'm using RabbitMQ to handle Pub/Sub communication inbetween services.

I'm using Google Kubernetes Engine for this project.

## Installation

You can try and run the project using docker-compose.yml included in the repo however it's not recommended approach since this will require substantial amount of computing power and memory from your local machine.

### Docker images
Docker images are available in my public repositories on Docker Hub.
You might as well build them yourself from Dockerfiles included in the repo, each located in a corresponding service directory.

### Kubernetes
To launch these services on your cluster you just need to clone the repo.

```
git clone https://github.com/Krixlion/microservices
```

And apply Kubernetes manifests located in K8s/ directory.
```
kubectl apply -R -f microservices/k8s/
```

You might need to tweek some settings to match your cloud provider. 
