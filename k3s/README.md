# Setting Up K3s Development Environment

## Prerequisites
* **Docker** installed
* **k3d** installed
* **kubectl** installed

## Creating K3s Cluster

### Create cluster with local registry and port mappings

```bash
k3d cluster create dev-cluster \
  --registry-create registry.localhost:5001 \
  --agents 1 \
  --port "50051:50051@loadbalancer" \
  --port "9090:9090@loadbalancer" \
  --port "9091:9090@loadbalancer"
```

## Configuring Kubeconfig

### Update kubeconfig file

```bash
k3d kubeconfig get dev-cluster > ~/.kube/k3d-dev.yaml &&
sed -i '' 's/host.docker.internal/127.0.0.1/g' ~/.kube/k3d-dev.yaml &&
export KUBECONFIG=~/.kube/k3d-dev.yaml
```

## Building and Pushing Docker Image

### Build and push image to local registry

```bash
docker build -t ton-node-service:latest . &&
docker tag ton-node-service:latest localhost:5001/ton-node-service:latest &&
docker push localhost:5001/ton-node-service:latest
```

## Creating Namespaces and ConfigMaps

### Create necessary namespaces

```bash
kubectl create namespace backend && kubectl create namespace monitoring
```

### Create config maps

```bash
kubectl create configmap ton-node-config --from-file=config.yaml=./config.yaml -n backend &&
kubectl create configmap prometheus-config --from-file=prometheus.yml=./infra/prometheus/prometheus-k3s.yml -n monitoring
```

## Deploying Services

### Apply deployment configurations

```bash
kubectl apply -f k3s/backend-deployment.yaml &&
kubectl apply -f k3s/backend-service.yaml &&
kubectl apply -f k3s/prometheus-deployment.yaml &&
kubectl apply -f k3s/prometheus-service.yaml
```

## Monitoring and Debugging

### Check pod status

```bash
kubectl get pods --all-namespaces
```

### Port forwarding

```bash
kubectl port-forward svc/ton-node-service -n backend 50051:50051 &
kubectl port-forward svc/ton-node-service -n backend 9090:9090 &
kubectl port-forward svc/ton-node-prometheus -n monitoring 9091:9090 &
```

## Useful Commands

**Basic Monitoring:**
```bash
kubectl get pods --all-namespaces
kubectl get ns
kubectl get all --all-namespaces
kubectl get pods -n <namespace>
kubectl get rs -n <namespace>
kubectl describe pod <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace>
```

**Deleting Resources:**
```bash
kubectl delete rs <pod-name> -n <namespace>
```

**Delete Pod After Update Config**
```bash
kubectl delete pod -n backend -l app=ton-node-service &&
kubectl delete pod -n backend -l app=ton-node-prometheus

kubectl apply -f k3s/backend-deployment.yaml &&
kubectl apply -f k3s/backend-service.yaml &&
kubectl apply -f k3s/prometheus-deployment.yaml &&
kubectl apply -f k3s/prometheus-service.yaml
```

## Rebuilding After Code Changes

**For Go Service Updates:**
1. Make changes in Go code
2. Rebuild Docker image:

```bash
docker build -t ton-node-service:latest . &&
docker tag ton-node-service:latest localhost:5001/ton-node-service:latest &&
docker push localhost:5001/ton-node-service:latest
```

3. Update deployment config with new image tag

```bash
kubectl rollout restart deployment/ton-node-service -n backend
```

4. Apply changes:

```bash
kubectl apply -f k3s/backend-deployment.yaml
```

## Complete Teardown

### Delete cluster and all resources
```bash
k3d cluster delete dev-cluster
docker ps -a | grep k3d | awk '{print $1}' | xargs docker rm -f
docker volume ls | grep k3d | awk '{print $2}' | xargs docker volume rm -f
docker rmi ton-node-service:latest
docker rmi localhost:5001/ton-node-service:latest
```
