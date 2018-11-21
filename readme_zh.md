# IAM - 账号管理和权限管理服务

IAM(Identity and Access Management)是账号管理和权限管理服务管理服务的简称. IAM模块对外通过Protobuf定义GRPC接口, 其它服务通过GRPC服务访问该模块的功能.

## 目录结构

- api: Protobuf规范文件, GRPC接口规范
- cmd: 每个服务的可执行程序
- doc: 文档和图片文件, 中文以 `_zh.md` 为后缀
- pkg: Go语言包文件

## IM - 账号管理服务

账号管理服务, 主要管理用户信息和组信息.

## AM - 权限管理服务

AM模块是基于RBAC(Role Based Access Control)模型提供权限管理服务.

![](./doc/images/rbac.dot.png)

RBAC重点涉及以下几个概念:

- 实体: 对应用户或者是类似用户的服务
- 资源: 对应请求的一个URL路径, 表示一个资源抽象, URL必须满足一定的规则
- 角色: 用于描述一类实体, 最终的权限是真的角色来设置的
- 角色绑定: 记录实体和角色的对应关系, 实体需要转为角色后才能被授权
- 规则: 规则描述一类资源的访问权限, 一组规则授权给角色

常见的OpenPitrix的API列表:

- `GET /api/v1/repos/repo-abcd/create_time`
- `GET /api/v1/users/user-name/runtimes`
- `GET /api/v1/runtimes/rt-abcd/cpu-num`

请求的数据流程
