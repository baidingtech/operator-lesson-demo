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

# Indexer原理
---


![width:26cm height:14cm](./images/design.png)

---
# 数据存储

`cache`委托`threadSafeMap`存放数据
```
type Indexer interface {
	Store
	Index(indexName string, obj interface{}) ([]interface{}, error)
	IndexKeys(indexName, indexedValue string) ([]string, error)
	ListIndexFuncValues(indexName string) []string
	ByIndex(indexName, indexedValue string) ([]interface{}, error)
	GetIndexers() Indexers
	AddIndexers(newIndexers Indexers) error
}
```

---
# 建立索引
```go
//key是IndexFunc计算出来的结果，比如default，value是所有obj的key的集合
type Index map[string]sets.String

//key是索引的分类名，比如namespace，value是一个方法，通过该方法可以获取obj的namespace，比如default
type Indexers map[string]IndexFunc

//key是索引的分类名，比如namespace
type Indices map[string]Index
```
![bg height:14cm right](./images/cache_and_indexer.png)

---

# 更新索引

通过`updateIndices`实现。

![bg height:14cm right](./images/indexupdate.png)

---


<!--
_class: lead
-->
### 谢谢