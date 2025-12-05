# gRPC 远程调用示例

本示例展示如何通过 gRPC 将 `Jinvic/pubsub` 发布订阅系统包装为远程服务，实现跨进程或跨网络的发布订阅功能。

## 项目结构

```bash
pubsub-grpc/
├── cmd/
│   ├── client/          # gRPC 客户端
│   │   ├── examples.go  # 示例代码
│   │   └── main.go      # 客户端入口
│   └── server/          # gRPC 服务端
│       └── main.go      # 服务端入口
├── pkg/
│   └── pb/              # Protocol Buffer 生成的 Go 代码
│       └── pubsub/
├── proto/               # Protocol Buffer 定义文件
│   └── pubsub.proto
├── go.mod
└── go.sum
```

## 功能特性

- **远程调用**：将本地 pubsub 系统通过 gRPC 暴露为远程服务
- **非阻塞消费**：提供 `TryConsume` 方法实现非阻塞消息消费
- **持续消费**：提供 `KeepConsume` 方法实现流式消息消费
- **完整生命周期**：支持发布者/订阅者的创建、订阅、发布、注销等完整操作

## RPC 服务接口

### 发布者相关

- `NewPublisher` - 创建新的发布者
- `UnregisterPublisher` - 注销发布者
- `Publish` - 发布消息到指定主题

### 订阅者相关

- `NewSubscriber` - 创建新的订阅者
- `UnregisterSubscriber` - 注销订阅者
- `Subscribe` - 订阅指定主题
- `UnsubscribePublisher` - 取消订阅整个发布者
- `UnsubscribeTopic` - 取消订阅指定主题

### 消息消费

- `TryConsume` - 尝试非阻塞消费消息
- `KeepConsume` - 持续消费消息（服务器流式）

## 运行示例

1. **启动服务端**：

   ```bash
   cd cmd/server
   go run main.go
   # 或指定端口
   go run main.go -port 12345
   ```

2. **运行客户端**：

   ```bash
   cd cmd/client
   go run main.go
   # 或连接到指定端口
   go run main.go -port 12345
   ```

## 示例说明

### 非阻塞消费示例

- 演示如何使用 `TryConsume` 方法进行非阻塞消费
- 首先尝试消费，由于没有消息返回空结果
- 发布消息后再次消费，成功获取消息

### 持续消费示例

- 演示如何使用 `KeepConsume` 方法持续接收消息
- 服务端通过流式响应持续发送消息

## 依赖

- [Jinvic/pubsub](https://github.com/Jinvic/pubsub) - 本地发布订阅系统
- [gRPC](https://grpc.io/) - 远程过程调用框架
- [Protocol Buffers](https://protobuf.dev/) - 序列化结构化数据
