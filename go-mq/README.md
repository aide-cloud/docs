# Go-MQ

一个轻量级的内存消息队列实现。

## 简介

Go-MQ 是一个用 Go 语言实现的轻量级内存消息队列库，主要用于进程内的消息传递和事件通知。它提供了一个简单但功能完整的消息队列接口，特别适合用于测试环境或者小型应用中的消息传递场景。


## 核心特性

- 简单的接口设计
- 基于 channel 的异步消息传递
- 支持多个主题（Topic）
- 支持优雅关闭
- 线程安全

## 快速开始

```go
// 创建 MQ 实例
mq := NewMockMQ()

// 订阅消息
msgChan := mq.Receive("test-topic")

// 发送消息
mq.Send("test-topic", []byte("Hello World"))

// 接收消息
msg := <-msgChan

// 清理资源
mq.RemoveReceiver("test-topic")
mq.Close()
```

```go
func TestNewMockMQ(t *testing.T) {
	mq := NewMockMQ()
	defer mq.Close()

	ch := mq.Receive("test")
	go func() {
		for msg := range ch {
			t.Logf("receive message: %s", msg.Data)
		}
		t.Log("receiver exit")
	}()
	for i := 0; i < 10; i++ {
		mq.Send("test", []byte("hello world "+strconv.Itoa(i)))
		time.Sleep(time.Second)
		if i == 5 {
			mq.RemoveReceiver("test")
		}
	}
}
```

## API 文档

### IMQ 接口

```go
type IMQ interface {
    // Send 发送消息
    Send(topic string, data []byte) error

    // Receive 接收消息 返回一个接收通道
    Receive(topic string) <-chan *Msg

    // RemoveReceiver 移除某个topic的接收通道
    RemoveReceiver(topic string)

    // Close 关闭连接
    Close()
}
```

### 消息结构

```go
type Msg struct {
    Data  []byte
    Topic []byte
}
```

## 适用场景

1. **单元测试**
   - 模拟消息队列行为
   - 测试消息处理逻辑

2. **原型开发**
   - 快速验证消息传递机制
   - 开发环境中使用

3. **小型应用**
   - 进程内的事件通知
   - 组件间的解耦

## 局限性

- 仅支持进程内的消息传递
- 消息不持久化
- 重启后消息丢失
- 不支持分布式场景

## 注意事项

1. 这是一个内存消息队列，所有数据都存储在内存中
2. 每个主题的默认缓冲区大小为 100
3. 在生产环境中建议使用成熟的消息队列系统（如 RabbitMQ 或 Kafka）
