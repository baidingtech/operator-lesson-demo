### Operator工具对比

前面我们对client-go/kubebuilder
进行了了解和学习，从这节课开始我们来了解另外一个编写Operator的工具operator-sdk
。它属于operator-framework的一个子项目，早期它与kubebuilder
是相互独立的，但底层又都是基于controller-runtime
进行的实现，后面为了社区更好的发展，两个社区进行了协作，现在我们看到的operator-sdk 
golang的实现基于kubebuilder来实现的。所以我们后面将着重了解operator-sdk
相对kubebuilder提供的独有功能，比如helm、ansible、olm。


本节课我们先来看看，我们可以使用哪些operator工具来帮我们实现operator。

![img.png](img.png)

其他Operator工具:
1. [Java Operator SDK](https://github.com/java-operator-sdk/java-operator-sdk)
2. [kudo](https://kudo.dev/)
3. [Kubernetes Operator Pythonic Framework (KOPF)](https://kopf.readthedocs.io/)
4. [Shell-Operator](https://github.com/flant/shell-operator)
5. [MetaController](https://github.com/metacontroller/metacontroller)
