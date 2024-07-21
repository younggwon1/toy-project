# Kubernetes-Resource-Manager

## The service is implemented only when the deployment status is error.

### 1. Down
#### 1.1 Enter the namespace for which you want to lower the replica count to 0.

```
go run main.go down --namespace {namespace}
```

EX)
```
go run main.go down --namespace "rnd-test"
```

#### 1.2 Enter the name of the deployment whose replica number you want to reduce to 0 and the namespace it belongs to.

```
go run main.go down --names "{deployment name one},{deployment name two},..."--namespace {namespace}
```

EX)
```
go run main.go down --names "test1" --namespace "rnd-test"
go run main.go down --names "test1,test2" --namespace "rnd-test"
```

### 2. Delete
#### 2.1 Enter the namespace for which you want to delete the replica 0 deployments.
The only situation in which a deployment can be deleted is if its replica is 0.

```
go run main.go delete --namespace {namespace}
```

EX)
```
go run main.go delete --namespace "rnd-test"
```