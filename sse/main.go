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
