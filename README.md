# DDD 分层架构 Demo

# 目录结构

```
.
├── application             # 应用层
│   └── service
├── common                  # 通用类库
│   ├── consts
│   │   └── consts.go
│   ├── converter
│   │   └── converter.go
│   ├── entity
│   │   ├── do
│   │   ├── dto
│   │   ├── po
│   │   ├── req
│   │   └── vo
│   ├── exceptions
│   └── response
├── document                # 文档
│   ├── docker
│   │   └── Dockerfile
│   └── swagger
├── domain                  # 领域层
│   └── service
├── infrastructure          # 基础服务层
│   ├── auth
│   ├── cache
│   ├── config
│   ├── mq
│   └── persistence
├── interface               # 用户接口层
│   └── web
│       └── gin
│           ├── controller
│           ├── middleware
│           └── router
├── config.yaml             # 配置文件
├── go.mod
└── main.go
├── README.md
```