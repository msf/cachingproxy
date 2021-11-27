// Package main implements the application initialization code
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/handler/mtcache"
	"github.com/msf/cachingproxy/handler/mtproxy"
	"github.com/msf/cachingproxy/server"
	ginlogrus "github.com/rocksolidlabs/gin-logrus"
	"github.com/sirupsen/logrus"
)

var (
	// Name injected at build time
	Name string

	// Version injected at build time
	Version string

	// BuildTime injected at build time
	BuildTime string

	// ListenPort
	ListenPort string

	// Logger
	log *logrus.Logger
)

func ServeHTTP() error {
	// TODO: adopt cobra
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		ginlogrus.Logger(log, "", true, true, os.Stdout, logrus.DebugLevel),
		gin.Recovery(),
	)
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithExcludedPaths([]string{"/metrics"})))

	// TODO: cmdline args for this
	srv, err := server.NewGinServer(
		mtcache.Config{
			MaxSizeMB: 20,
			MaxTTL:    72 * time.Hour,
		},
		map[mtproxy.RoutingKey]string{
			{SourceLang: "en", TargetLang: "pt"}: "bananas.foo",
			{}:                                   "bar.foo",
		},
	)
	if err != nil {
		return err
	}

	r.GET("/ping", srv.Ping)
	r.GET("/echo/:id/:cnt", srv.Message)
	r.POST("/v1/machine_translate", srv.MachineTranslate)

	return r.Run(fmt.Sprintf(":%s", ListenPort))
}

// init is only used for keeping the command setup within the same file
func main() {
	ListenPort = "4321"
	log = logrus.New()

	log.WithFields(logrus.Fields{
		"Name":       Name,
		"Version":    Version,
		"BuildTime":  BuildTime,
		"ListenPort": ListenPort,
	}).Print("Starting now")

	if err := ServeHTTP(); err != nil {
		log.Error("ServeHTTP error", err)
	}
}
