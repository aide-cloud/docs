package main

import (
	"fmt"
	"log"
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

	// 跨域
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

	r.GET("/say", func(ctx *gin.Context) {
		msg := msgFN()
		msgCh <- msg
		ctx.JSON(200, map[string]string{
			"message": msg,
		})
	})

	// SSE 路由
	r.GET("/sse", func(c *gin.Context) {
		// 设置响应头
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// 模拟 SSE 数据推送
		for {
			// 客户端断开链接后退出
			select {
			case <-c.Request.Context().Done():
				log.Println("client disconnected")
				return
			case msg := <-msgCh:
				fmt.Fprintf(c.Writer, "data: %s\n\n", msg)
				c.Writer.Flush()
			default:
				continue
			}
		}
	})

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
