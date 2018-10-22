# webgo 

Webgo 是一个微内核的Go语言Web框架，是对网上go开发中在总结，

模块化的处理每一个任务和对象，使用go语言中间件模式来完成每次氢气的任务

## 功能模块

recovery - 可以是的代码在panic后继续执行
logger - 模仿python的logging库实现一个日志系统
context - context保存在一个全局的map类型中来管理上下文
session - 用来处理数据库链接的，将一个session绑定到当前请求上
xsrf - Generates and validates csrf tokens
auth - Build Status Coverage Status manage permissions via ACL, RBAC, ABAC


## todo

1. 缓存: http, redis
2. 上下文
