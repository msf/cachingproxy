// Package main implements the application initialization code
package main

import (
	"strings"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e := echo.New()
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 1,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/metrics")
		},
	}))

	e.GET("/echo/:id/:cnt", server.EchoMessage)
	e.GET("/ping", server.EchoPing)
	return e.Start(":1323")
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
