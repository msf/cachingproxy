package cmd

import "github.com/gin-gonic/gin"

const StatusOK = 200

func ServeHTTP() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r.Run(":4321")
}
