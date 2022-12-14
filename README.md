# DDD 分层架构 Demo

## 目录结构

```
.
├── application                 # 应用层
│   └── service
├── common                      # 通用类库
│   ├── consts
│   │   └── consts.go
│   ├── converter
│   │   └── converter.go
│   ├── encrypt
│   │   └── aes.go
│   ├── entity
│   │   ├── do
│   │   ├── dto
│   │   ├── po
│   │   ├── req
│   │   └── vo
│   ├── factory
│   │   └── object_factory.go
│   ├── fsm
│   │   ├── errors.go
│   │   ├── event.go
│   │   └── fsm.go
│   ├── response
│   │   └── response.go
│   ├── serializer
│   │   └── encode.go
│   └── utils.go
├── document                    # 文档
│   ├── docker
│   │   └── Dockerfile
│   └── swagger
│       └── doc
├── domain                      # 领域层
│   └── service
├── infrastructure              # 基础服务层
│   ├── auth
│   │   └── auth.go
│   ├── cache
│   │   ├── big_cache.go
│   │   └── cache.go
│   ├── config
│   │   ├── config.go
│   │   └── yaml_config.go
│   ├── es
│   │   └── es_client.go
│   ├── http
│   │   └── resty_client.go
│   ├── mq
│   │   └── kafka_client.go
│   ├── persistence
│   │   ├── mongodb_client.go
│   │   ├── mysql_client.go
│   │   └── redis_client.go
│   └── singleflight
│       └── single_flight.go
├── interface                   # 用户接口层
│   └── web
│       └── gin
│           ├── controller
│           ├── middleware
│           └── router
├── main.go                     # 主函数
├── config.yaml                 # 配置文件
├── go.mod
├── go.sum
└── README.md
```