package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Echo(c *gin.Context) {
	// TODO: pass data from gin.Context
	m := model.Message{
		ID:      "id",
		Content: "content",
	}
	r, err := handler.EchoMessage(m)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// TODO: convert reply to Gin response
	c.JSON(http.StatusOK, r)
}
