package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	_ = router.Run("127.0.0.1:8082") // 监听并在 0.0.0.0:8080 上启动服务
}
