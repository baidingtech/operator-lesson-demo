# 用Operator-SDK：Helm实现Operator

### 创建helm chart

1. 我们先将app的应用创建一个helm chart。

```shell
helm create app-chart
```
> helm chart可以放在本地或者远程。
 
2. 查看chart的结果是否正确

```shell
helm template app-chart
```


### 安装operator-sdk

参考https://sdk.operatorframework.io/docs/installation/

### 创建项目

> 可通过 operator-sdk init --plugins helm --help 查看完整使用方法

```shell
mkdir helm
cd helm
operator-sdk init --plugins helm --domain baiding.tech --group ingress --version v1alpha1 --kind App --helm-chart ../app-chart
```

### 测试operator

1. 安装crd并在本地启动operator

```shell
make install
make run
```
2. 修改config/samples/ingress_v1alpha1_app.yaml进行验证

### 部署operator

1. 构建镜像和部署到集群

```shell
make docker-build IMG=wangtaotao2015/operator-helm-demo
make docker-push IMG=wangtaotao2015/operator-helm-demo
make deploy IMG=wangtaotao2015/operator-helm-demo
```

2. 修改config/samples/ingress_v1alpha1_app.yaml进行验证









### 参考文档

https://sdk.operatorframework.io/docs/building-operators/helm/