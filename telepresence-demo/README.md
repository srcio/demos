# telepresence

[Telepresence 官方文档](https://www.telepresence.io/docs/latest/quick-start/)

## 简介

`telepresence` 是一个开源工具，它可以让你在本地运行单个服务，同时将该服务连接到远程 `Kubernetes` 集群。这让开发人员可以在多服务应用程序上工作:

-   对单个服务进行快速本地开发，即使该服务依赖于集群中的其他服务。对服务进行更改并保存，您可以立即看到新服务的运行情况。

-   使用任何安装在本地的工具来测试/调试/编辑您的服务。例如，您可以使用调试器或IDE!

-   获得本地开发机器就像运行 `Kubernetes` 集群中一样的体验：接管集群中上游微服务组件的请求，直接使用集群 `service` 名称网络访问下游微服务组件；

## 工作原理

`telepresence` 在 `Kubernetes` 集群中部署了一个双向网络代理服务。这个服务将 `Kubernetes` 环境中的数据（例如 TCP 连接、环境变量、卷）代理到本地进程。本地进程透明地覆盖其网络，以便DNS调用和TCP连接通过代理路由到远程Kubernetes集群。

## 安装

**Windows**

powershell 运行命令：

```powershell
curl -fL https://app.getambassador.io/download/tel2/windows/amd64/latest/telepresence.zip -o telepresence.zip
Expand-Archive -Path telepresence.zip
Remove-Item 'telepresence.zip'
cd telepresence
Set-ExecutionPolicy Bypass -Scope Process
.\install-telepresence.ps1
cd ..
Remove-Item telepresence
```

**Linux**

```bash
sudo curl -fL https://app.getambassador.io/download/tel2/linux/amd64/latest/telepresence -o /usr/local/bin/telepresence
sudo chmod a+x /usr/local/bin/telepresence
```

**MacOS**

```bash
# Intel Macs
brew install datawire/blackbird/telepresence

# Apple silicon Macs
brew install datawire/blackbird/telepresence-arm64
```

## 调试

![image-20220810112359349](https://srcio.oss-cn-hangzhou.aliyuncs.com/images/image-20220810112359349.png)

假设我们有如上图中一个微服务运行在 K8s 集群中：

`service-a => service-b => service-c`（图中 `local-service-b` 为开发机待调试服务）

使用如下命令将这个模拟微服务运行起来：

```bash
curl -s https://raw.githubusercontent.com/srcio/telepresence-demo/master/deploy.sh | bash
```
我们需要调试微服务其中的一个组件 `service-b`，那么 `telepresence` 可以成为你的得力助手。

1.   **在本地开发机使用 `telepresence` 连接到集群**

```bash
telepresence connnect
```

>   本地需要可以访问集群的 kubeconfig 配置文件：`~/.kube/config`，或者使用 `--kubeconfig` 指定配置文件。

2.   **运行本地待调试应用 `service-b`**

```
go run service-b/main.go
```
>   该服务端口默认使用 80。

3.   **代理集群中目标微服务组件 `service-b` 的网络**

```bash
telepresence intercept service-b --service service-b --port 80:80
```

> 其中：
>
> `--service` 参数后面值代表集群中目标微服务组件的 `Service` 资源名称；
>
> `--port` 参数后面的两个端口分别代表：{本地调试服务端口}:{集群微服务组件 `Service` 端口}。

4.   **调试**

由于使用了 `telepresence` 的连接，我们可以直接在本地连接 `service-a` 服务：

```bash
curl service-a.default
```

我们知道 `service-a => service-b => service-c` 调用关系，运行以上命令后观察本地调试应用，可以看到如下日志：

```tex
service-b started...
request from service-a
Response from service-c! 2022-08-10 03:17:37.7909436 +0000 UTC m=+550.818222118
```

说明我们使用 `telepresence` 的本地调试之成功了！

5.   **思考**

你可能会有疑问，现在 `service-a` 的请求转发到了本地调试应用，那么集群中 service-b 还会接受到请求吗？

如果你细心的查看一下集群中 `service-b` 的 `Pod`，你会发现 `Pod` 中多出了一个 `traffic-agent` 容器，以及一个 `tel-agent-init` 的初始化容器，这两个多出的容器就是帮助接管 `service-b` 的微服务容器的请求的；

再查看一下 `service-b` 的微服务容器日志，发现里面并不会接收到来自 `service-a` 的请求。

6.   **回到解放前**

```
telepresence leave service-b
telepresence quit
curl -s https://raw.githubusercontent.com/srcio/telepresence-demo/master/clean.sh | bash
```

## RBAC

[Telepresence RBAC 文档](https://www.telepresence.io/docs/latest/reference/rbac/)

开发人员在使用 telepresence 在本地开发和调试服务，应该遵循最小权限原则。

如下的 rbac 配置，可以限制开发人员直允许对 `test` 命名空间下组件进行调试，并且不具备对 Pod 以外的资源的修改删除权限。

*rbac/tp-namespaced-rbac-user.yaml*

```bash
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tp-user                                       # Update value for appropriate user name
  namespace: ambassador                                # Traffic-Manager is deployed to Ambassador namespace
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name:  telepresence-role
rules:
- apiGroups:
  - ""
  resources: ["pods"]
  verbs: ["get", "list", "create", "watch", "delete"]
- apiGroups:
  - ""
  resources: ["services"]
  verbs: ["update"]
- apiGroups:
  - ""
  resources: ["pods/portforward"]
  verbs: ["create"]
- apiGroups:
  - "apps"
  resources: ["deployments", "replicasets", "statefulsets"]
  verbs: ["get", "list", "update", "watch"]
- apiGroups:
  - "getambassador.io"
  resources: ["hosts", "mappings"]
  verbs: ["*"]
- apiGroups:
  - ""
  resources: ["endpoints"]
  verbs: ["get", "list", "watch"]
---
kind: RoleBinding                                      # RBAC to access ambassador namespace
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: t2-ambassador-binding
  namespace: ambassador
subjects:
- kind: ServiceAccount
  name: tp-user                                       # Should be the same as metadata.name of above ServiceAccount
  namespace: ambassador
roleRef:
  kind: ClusterRole
  name: telepresence-role
  apiGroup: rbac.authorization.k8s.io
---
kind: RoleBinding                                      # RoleBinding T2 namespace to be intecpeted
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: telepresence-test-binding                      # Update "test" for appropriate namespace to be intercepted
  namespace: test                                      # Update "test" for appropriate namespace to be intercepted
subjects:
- kind: ServiceAccount
  name: tp-user                                       # Should be the same as metadata.name of above ServiceAccount
  namespace: ambassador
roleRef:
  kind: ClusterRole
  name: telepresence-role
  apiGroup: rbac.authorization.k8s.io
​
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name:  telepresence-namespace-role
rules:
- apiGroups:
  - ""
  resources: ["namespaces"]
  verbs: ["get", "list", "watch"]
- apiGroups:
  - ""
  resources: ["services"]
  verbs: ["get", "list", "watch"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: telepresence-namespace-binding
subjects:
- kind: ServiceAccount
  name: tp-user                                       # Should be the same as metadata.name of above ServiceAccount
  namespace: ambassador
roleRef:
  kind: ClusterRole
  name: telepresence-namespace-role
  apiGroup: rbac.authorization.k8s.io
```

由于 telepresence 连接集群时需要 kubeconfig 文件，所以我们需要根据当前开发人员所分配的 ServiceAccount 生成 kubeconfig 文件：

```bash
kubectl apply -f ./rbac/tp-namespaced-rbac-user.yaml
./k2 config gen --from tp-user -n ambassador > tp_kubeconfig

telepresence connect --kubeconfig tp_kubeconfig
```

