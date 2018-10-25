# webgo 

web.go 是一个微内核的Go语言RESTfil框架，可以实现基本的API设计功能

模块化的处理每一个任务和对象，使用go语言中间件模式来完成每次氢气的任务

## 路由设计

路由设计中将每次的资源和uri一一对应之后再保存在一个全局的map中

```
map[regexp.Regexp]Resource

type Resource interface {
    Match() bool
}

func(rs *Resource) get(w http.ResponseWriter, r *http.Request) {
    ...
}
```
这样做的好处是:

1. 通过正则匹配路由来实现动态路由和路由参数设置
2. 我们可以通过反射来对应每个请求方法, 在资源中每个方法对应一个handler

示例:

```
/api/users, users
/api/users/:id, users
/api/users/:id/books, books
/api/books, books
/api/books/:id, books
```

实现:

```
router.add(/api/users/(\w), users)
```


## 功能模块

recovery - 可以是的代码在panic后继续执行
logger - 模仿python的logging库实现一个日志系统
context - context保存在一个全局的map类型中来管理上下文
session - 用来处理数据库链接的，将一个session绑定到当前请求上
xsrf - Generates and validates csrf tokens
auth - Build Status Coverage Status manage permissions via ACL, RBAC, ABAC


## todo


1. 动态路由
2. 缓存: http, redis
3. 上下文
