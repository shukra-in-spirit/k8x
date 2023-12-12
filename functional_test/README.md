## Functional Testing
Minimum go version required is 1.19, other tool required is kubectl, minikube

### Infrastructure
1. Run the following commands to deploy the prometheus community stack in default namespace:
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install [RELEASE_NAME] prometheus-community/kube-prometheus-stack
kubectl get pods
```
2. There would be pods which are showing error in pulling images. The images are hosted in the quay.io repository, it gives issues with minikube. we can manually download the images and load them.
```
kubectl desribe pod [POD_NAME]
# check the image name at the bottom
docker pull [IMAGE_NAME]
minikube image load [IMAGE_NAME]
```
3. Repeat until all pods are in running state.

### Tests

Please follow the following steps to test:

1. Run the following commands from the root of the repository.
```
kubectl create ns test-namespace
kubectl apply -f functional_test/test-deployment.yaml
kubectl get pods -n test-namespace
```
2. Run the command `go run functional_test/dummy_main.go`.

### Functional Tests
1. The scaling interface created on top of the golang kubernetes client has been tested. All the methods in the interface are called and their response time logged in milliseconds.

