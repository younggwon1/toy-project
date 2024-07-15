# gitops golang
> This is a service implemented so that multiple deployment pipelines can be easily deployed using one tool.

### Deploy to Kubernetes
#### Required env
1. GIT_USERNAME
2. GIT_PASSWORD
3. ARGOCD_SERVER
4. ARGOCD_TOKEN

#### Required Template for Deploying to Kubernetes
```
spec:
  kubernetes:
    deploys:
    - helm:
        url: "helm-url-1"
        org: "helm-org-1"
        repo: "helm-repo-1"
        values:
        - file: "value-file-1"
        - file: "value-file-2"
      argocd:
        url: "argocd-url"
        apps:
        - name: "app-name-1"
        - name: "app-name-2"
```

#### gitops golang cli for deploying to kubernetes
```
go run main.go deploy \
    --user "${git user}" \
    --email "${git email}" \
    --values image values \
    --spec spec path
```

example
```
go run main.go deploy \
    --user "younggwon" \
    --email "younggwon@aaa.bbb" \
    --values "{\"image.tag\":\"dev-12345\"}" \
    --spec "../../deploy-dev.yaml"
```




go run main.go deploy \
    --spec "../../deploy-dev.yaml"