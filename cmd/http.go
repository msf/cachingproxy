package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/server"
)

func ServeHTTP() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", server.Ping)
	return r.Run(":4321")
}
