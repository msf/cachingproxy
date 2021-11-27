package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/handler"
	mt "github.com/msf/cachingproxy/handler/mt"
	"github.com/msf/cachingproxy/handler/mtcache"
	"github.com/msf/cachingproxy/handler/mtproxy"
	"github.com/msf/cachingproxy/model"
)

type GinServer struct {
	mtHandler handler.MachineTranslationHandler
}

func NewGinServer(cacheConfig mtcache.Config,
	routingMap map[mtproxy.RoutingKey]string,
) (*GinServer, error) {
	mtH, err := mt.NewCachingMTHandler(cacheConfig, routingMap)
	if err != nil {
		return nil, err
	}
	return &GinServer{
		mtHandler: mtH,
	}, nil
}

func (s *GinServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s *GinServer) Message(c *gin.Context) {
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

func (s *GinServer) MachineTranslate(c *gin.Context) {
	var req model.MachineTranslationRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	r, err := s.mtHandler.Handle(&req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, r)
}
