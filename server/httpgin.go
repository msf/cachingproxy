package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
)

func GinPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GinMessage(c *gin.Context) {
	m := model.Message{
		ID:      c.Param("id"),
		Content: c.Param("cnt"),
	}
	r, err := handler.EchoMessage(m)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, r)
}

func GinMTRoute(c *gin.Context) {
	var req model.MachineTranslationRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	r, err := handler.MachineTranslate(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, r)
}
