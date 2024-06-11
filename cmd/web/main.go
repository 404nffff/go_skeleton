package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
)

func main() {
	// 创建一个新的 Gin 引擎
	r := gin.Default()

	// 创建一个 Ants 池，错误检查
	pool, err := ants.NewPool(10)
	if err != nil {
		log.Fatalf("Failed to create ants pool: %v", err)
	}
	defer pool.Release()

	// 使用中间件管理并发和错误处理
	r.Use(func(c *gin.Context) {
		// 创建一个 channel 用于接收 goroutine 的完成信号
		done := make(chan struct{})
		panicChan := make(chan interface{})

		go func() {
			defer close(done)
			defer func() {
				if r := recover(); r != nil {
					panicChan <- r
				}
			}()
			c.Next()
		}()

		select {
		case <-done:
			// 请求正常完成
		case p := <-panicChan:
			log.Printf("Recovered from panic: %v", p)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	})

	// 定义一个处理函数
	handler := func(c *gin.Context) {
		// 获取请求参数
		param := c.Query("param")
		if param == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "param is required"})
			return
		}

		// 从池中获取一个 goroutine
		err := pool.Submit(func() {
			// 调用你的任务方法
			result := myTask(param)

			// 使用 c.Request.Context() 确保上下文未被取消
			if c.Request.Context().Err() == nil {
				c.JSON(http.StatusOK, gin.H{"result": result})
			} else {
				log.Printf("Request context was cancelled")
			}
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to submit task"})
		}
	}

	// 路由
	r.GET("/task", handler)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// 定义你的任务方法
func myTask(param string) string {
	// 这里可以放你的任务处理逻辑
	return fmt.Sprintf("Processed: %s", param)
}
