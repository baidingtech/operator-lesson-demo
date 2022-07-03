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
# controller-tools介绍

---

### 安装

1. 获取controller-tools的代码，并切换到[v0.8.0 ](https://github.com/kubernetes-sigs/controller-tools/releases/tag/v0.8.0)的tag上

   ```shell
   git checkout v0.8.0
   ```

2. 编译项目，安装代码生成工具，这里我们只安装我们接下来会用到的工具

   ```shell
   go install ./cmd/{controller-gen,type-scaffold}
   ```

---
### 使用type-scaffold

```
type-scaffold --kind Foo
type-scaffold --help
```

---

### 使用controller-gen

根据标记生成代码

```
controller-gen crd paths=./... output:crd:dir=config/crds 
controller-gen object paths=./...
```

- https://book.kubebuilder.io/reference/controller-gen.html
- https://book.kubebuilder.io/reference/generating-crd.html
- [markers](https://book.kubebuilder.io/reference/markers.html)

---

     
### demo
<!--
_class: lead
-->

show me the code

---
<!--
_class: lead
-->

### 谢谢
