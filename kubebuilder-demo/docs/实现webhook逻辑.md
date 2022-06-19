# 使用kubebuilder实现webhook

### 需求

上节课我们实现了App的控制器的逻辑，但是我们希望在用户创建App资源时，
做一些更细的控制。


```yaml
apiVersion: ingress.baiding.tech/v1beta1
kind: App
metadata:
  name: app-sample
spec:
  image: nginx:latest
  replicas: 3
  enable_ingress: false #默认值为false，需求为：设置为反向值;为true时，enable_service必须为true
  enable_service: false
```

简单的校验我们可以直接使用CRD的scheme校验，但是复杂一点的需求我们
应该如何处理呢？
```yaml
//+kubebuilder:default:enable_ingress=false
```

我们这节课将会通过K8S提供的准入控制来实现。

### 准入控制

前面我们有学习到k8s提供了一系列的准入控制器，通过它们我们可以对api server的请求
进行处理。而对于我们自定义的需求，可以通过[MutatingAdmissionWebhook](https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook)
和[ValidatingAdmissionWebhook](https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook)
进行处理。

kubebuilder对其进行了支持，我们可以很方便的通过其实现我们的webhook逻辑。

### 创建webhook
> 创建webhook之前需要先创建api

 
1. 生成代码

```shell
kubebuilder create webhook --group ingress --version v1beta1 --kind App --defaulting --programmatic-validation
```

创建之后，在main.go中会添加以下代码:

```go
	if err = (&ingressv1beta1.App{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "App")
		os.Exit(1)
	}
```

同时会生成下列文件，主要有：

- api/v1beta1/app_webhook.go webhook对应的handler，我们添加业务逻辑的地方
  
- api/v1beta1/webhook_suite_test.go 测试
  
- config/certmanager 自动生成自签名的证书，用于webhook server提供https服务

- config/webhook 用于注册webhook到k8s中

- config/crd/patches 为conversion自动注入caBoundle

- config/default/manager_webhook_patch.yaml 让manager的deployment支持webhook请求
- config/default/webhookcainjection_patch.yaml 为webhook server注入caBoundle

注入caBoundle由cert-manager的[ca-injector](https://cert-manager.io/docs/concepts/ca-injector/#examples) 组件实现

2. 修改配置

为了支持webhook，我们需要修改config/default/kustomization.yaml将相应的配置打开，具体可参考注释。
```yaml
# Adds namespace to all resources.
namespace: kubebuilder-demo-system

# Value of this field is prepended to the
# names of all resources, e.g. a deployment named
# "wordpress" becomes "alices-wordpress".
# Note that it should also match with the prefix (text before '-') of the namespace
# field above.
namePrefix: kubebuilder-demo-

# Labels to add to all resources and selectors.
#commonLabels:
#  someName: someValue

bases:
- ../crd
- ../rbac
- ../manager
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
- ../webhook
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
- ../certmanager
# [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'.
#- ../prometheus

patchesStrategicMerge:
# Protect the /metrics endpoint by putting it behind auth.
# If you want your controller-manager to expose the /metrics
# endpoint w/o any authn/z, please comment the following line.
- manager_auth_proxy_patch.yaml

# Mount the controller config file for loading manager configurations
# through a ComponentConfig type
#- manager_config_patch.yaml

# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in
# crd/kustomization.yaml
- manager_webhook_patch.yaml

# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
# Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
# 'CERTMANAGER' needs to be enabled to use ca injection
- webhookcainjection_patch.yaml

# the following config is for teaching kustomize how to do var substitution
vars:
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
- name: CERTIFICATE_NAMESPACE # namespace of the certificate CR
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1
    name: serving-cert # this name should match the one in certificate.yaml
  fieldref:
    fieldpath: metadata.namespace
- name: CERTIFICATE_NAME
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1
    name: serving-cert # this name should match the one in certificate.yaml
- name: SERVICE_NAMESPACE # namespace of the service
  objref:
    kind: Service
    version: v1
    name: webhook-service
  fieldref:
    fieldpath: metadata.namespace
- name: SERVICE_NAME
  objref:
    kind: Service
    version: v1
    name: webhook-service

```

### webhook业务逻辑

#### 设置enable_ingress的默认值
```go
func (r *App) Default() {
	applog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	r.Spec.EnableIngress = !r.Spec.EnableIngress
}
```

#### 校验enable_service的值

```go
func (r *App) ValidateCreate() error {
	applog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateApp()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *App) ValidateUpdate(old runtime.Object) error {
	applog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateApp()
}

func (r *App) validateApp() error {
	if !r.Spec.EnableService && r.Spec.EnableIngress {
		return apierrors.NewInvalid(GroupVersion.WithKind("App").GroupKind(), r.Name,
			field.ErrorList{
				field.Invalid(field.NewPath("enable_service"),
					r.Spec.EnableService,
					"enable_service should be true when enable_ingress is true"),
			},
		)
	}
	return nil
}
```
### 测试

1. 安装cert-manager

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
```

2. 部署

```shell
IMG=wangtaotao2015/app-controller make deploy
```

3. 验证

```yaml
apiVersion: ingress.baiding.tech/v1beta1
kind: App
metadata:
  name: app-sample
spec:
  image: nginx:latest
  replicas: 3
  enable_ingress: false #会被修改为true
  enable_service: false #将会失败

```

```yaml
apiVersion: ingress.baiding.tech/v1beta1
kind: App
metadata:
  name: app-sample
spec:
  image: nginx:latest
  replicas: 3
  enable_ingress: false #会被修改为true
  enable_service: true #成功

```

```yaml
apiVersion: ingress.baiding.tech/v1beta1
kind: App
metadata:
  name: app-sample
spec:
  image: nginx:latest
  replicas: 3
  enable_ingress: true #会被修改为false
  enable_service: false #成功

```

### 如何本地测试

1. 添加本地测试相关的代码

- config/dev

- Makefile

```shell
.PHONY: dev
dev: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/dev | kubectl apply -f -
.PHONY: undev
undev: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/dev | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

```


2. 获取证书放到临时文件目录下

```shell
kubectl get secrets webhook-server-cert -n  kubebuilder-demo-system -o jsonpath='{..tls\.crt}' |base64 -d > certs/tls.crt
kubectl get secrets webhook-server-cert -n  kubebuilder-demo-system -o jsonpath='{..tls\.key}' |base64 -d > certs/tls.key
```

3. 修改main.go,让webhook server使用指定证书

```go
	if os.Getenv("ENVIRONMENT") == "DEV" {
		path, err := os.Getwd()
		if err != nil {
			setupLog.Error(err, "unable to get work dir")
			os.Exit(1)
		}
		options.CertDir = path + "/certs"
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
```

4. 部署

```shell
make dev
```

5. 清理环境

```shell
make undev
```



