// Package main implements the application initialization code
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/server"
	log "github.com/sirupsen/logrus"
)

var (
	// Name injected at build time
	Name string

	// Version injected at build time
	Version string

	// BuildTime injected at build time
	BuildTime string
)

func ServeHTTP() error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", server.GinPing)
	r.GET("/echo/:id/:cnt", server.GinMessage)
	return r.Run(":4321")
}

// init is only used for keeping the command setup within the same file
func main() {
	log.WithFields(log.Fields{
		"Name":      Name,
		"Version":   Version,
		"BuildTime": BuildTime,
	}).Print("Starting now")

	if err := ServeHTTP(); err != nil {
		log.Error("ServeHTTP error", err)
	}
}
