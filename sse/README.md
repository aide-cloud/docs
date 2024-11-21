# 使用 Go 和 HTML5 实现 SSE 实时推送功能

## 简介

Server-Sent Events (SSE) 是一种服务器推送技术，允许服务器实时地向浏览器推送数据。与 WebSocket 不同，SSE 是单向的（只能从服务器向客户端推送），但实现更简单，特别适合于实时通知、实时日志等场景。

## 技术栈

- 后端：Go + Gin 框架
- 前端：HTML5 + JavaScript
- 通信：SSE (Server-Sent Events)

## 功能特点

- 服务器实时推送消息到浏览器
- 自动重连机制
- 跨域支持 (CORS)
- 轻量级实现

## 代码实现

### 后端实现 (Go)

后端使用 Gin 框架实现了一个简单的 SSE 服务器：

- 创建消息通道用于数据传输
- 实现 CORS 中间件支持跨域访问
- 设置 SSE 相关的 HTTP 头
- 实现消息推送逻辑

```go
// sse/main.go
package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var msgCh = make(chan string, 100)
var msgFN = func() func() string {
	i := 0
	return func() string {
		i++
		return "hello " + strconv.Itoa(i)
	}
}()

func main() {
	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// SSE 端点
	r.GET("/sse", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		c.Writer.Flush()

		for {
			select {
			case <-c.Request.Context().Done():
				return
			case msg := <-msgCh:
				fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
				c.Writer.Flush()
			}
		}
	})

	r.Run(":8080")
}

```

### 前端实现 (HTML5)

前端使用原生的 EventSource API 实现 SSE 客户端：

- 建立 SSE 连接
- 处理接收到的消息
- 错误处理机制

```html
<!-- sse/index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SSE Demo</title>
</head>
<body>
<h1>Server-Sent Events</h1>
<div id="messages"></div>

<script>
    const eventSource = new EventSource("http://localhost:8080/sse");

    eventSource.onmessage = function (event) {
        const messagesDiv = document.getElementById("messages");
        const newMessage = document.createElement("div");
        newMessage.textContent = event.data;
        messagesDiv.appendChild(newMessage);
    };

    eventSource.onerror = function () {
        console.error("SSE connection error.");
    };
</script>
</body>
</html> 
```

## 运行说明

1. 启动后端服务：

```bash
cd sse
go run main.go
```

2. 打开前端页面：
   - 直接在浏览器中打开 `index.html` 文件
   - 或者使用简单的 HTTP 服务器托管

3. 发送消息

```bash
curl http://localhost:8080/say
```

## 实现原理

1. SSE 连接建立：
   - 客户端通过 EventSource API 连接到服务器的 /sse 端点
   - 服务器设置相应的 HTTP 头，建立持久连接

2. 消息推送：
   - 服务器通过 channel 接收要推送的消息
   - 使用特定格式（data: message\n\n）推送给客户端
   - 客户端通过 onmessage 事件处理接收到的消息

3. 连接管理：
   - 服务器监控连接状态，在客户端断开时及时清理资源
   - 客户端支持自动重连机制

## 注意事项

1. 浏览器兼容性：
   - EventSource API 在现代浏览器中得到广泛支持
   - IE 不支持 SSE，需要使用 polyfill 或替代方案

2. 性能考虑：
   - SSE 会占用服务器连接
   - 建议设置合适的超时和重连策略
   - 考虑使用负载均衡处理大量连接

## 扩展建议

1. 添加心跳机制确保连接存活
2. 实现消息重试机制
3. 添加消息 ID 和事件类型
4. 实现消息过滤和订阅机制
5. 添加安全认证机制

## 参考资料

- [MDN - Server-Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [HTML5 SSE 规范](https://html.spec.whatwg.org/multipage/server-sent-events.html)
