---
marp: true
theme: gaia
paginate: true
footer: '@白丁云原生'
backgroundColor: white
style: |
    code {
        background: black;
    }
---
<!--
_class: lead
-->

# ShareInformer原理
---

![width:26cm height:14cm](./images/design.png)

---

# ShareInformer的作用

主要负责完成两大类功能:
1. 缓存我们关注的资源对象的最新状态的数据
eg. 创建Indexer/Clientset(通过listerwatcher)/DeltaFIFO/Controller(包含Reflector的创建)

2. 根据资源对象的变化事件来通知我们注册的事件处理方法
eg. 创建sharedProcessor/注册事件处理方法

---

# ShareInformer的创建

- NewSharedIndexInformer
创建Informer的基本方法
- NewDeploymentInformer
创建内建资源对象对应的Informer的方法，调用NewSharedIndexInformer实现

- NewSharedInformerFactory
工厂方法，内部有一个map存放我们创建过的Informer，达到共享informer的目的，避免重复创建informer对象，浪费内存 

---

# ShareInformer的使用

```
	//create config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	//...
	//create client
	clientset, err := kubernetes.NewForConfig(config)
	//...
	//create informer
	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Pods().Informer()

	//register event handler
	informer.AddEventHandler()
	//...
	//start factory
	factory.Start(stopCh)
```

---


<!--
_class: lead
-->
### 谢谢