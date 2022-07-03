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
# api和apimachinery介绍

---

### api

https://github.com/kubernetes/api

主要功能：
- 内建资源对象定义
- 内建资源对象注册

思考： 自定义如何注册?

---
### apimachinery

https://github.com/kubernetes/apimachinery

主要存放服务端和客户端公用库，包含：
- ObjectMeta与TypeMeta
- Scheme
- RESTMapper
- 编码与解码
- 版本转换
- ...

---
### apimachinery: ObjectMeta与TypeMeta

- 与type对象的关系
- 与Object接口的关系


---
### apimachinery: Scheme

- type对象注册
- type对象与GVK的转换
- 默认值处理方法注册
- 版本转换方法注册

---

### apimachinery: RESTMapper

- GVK与GVR转换


---
<!--
_class: lead
-->

### 谢谢
