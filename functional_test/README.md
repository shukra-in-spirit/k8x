## Functional Testing
Minimum go version required is 1.19, other tool required is kubectl, minikube

Please follow the following steps to test:

1. Run the following commands from the root of the repository.
```
kubectl create ns test-namespace
kubectl apply -f functional_test/test-deployment.yaml
kubectl get pods -n test-namespace
```
2. Run the command `go run functional_test/dummy_main.go`.

### Tests

1. The scaling interface created on top of the golang kubernetes client has been tested. All the methods in the interface are called and their response time logged in milliseconds.
