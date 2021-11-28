package cmd

import (
	"fmt"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/msf/cachingproxy/handler/mtcache"
	"github.com/msf/cachingproxy/handler/mtproxy"
	"github.com/msf/cachingproxy/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	ginlogrus "github.com/toorop/gin-logrus"
)

var (
	GinPort     int16
	ReleaseMode bool
)

func init() {
	rootCmd.AddCommand(ginCmd)
	ginCmd.Flags().Int16Var(&GinPort, "ginPort", 4321, "gin server listening port")
	ginCmd.Flags().BoolVar(&ReleaseMode, "release", false, "release mode")
}

var ginCmd = &cobra.Command{
	Use:   "gin",
	Short: "using gin server",
	Run: func(cmd *cobra.Command, args []string) {

		log.WithFields(logrus.Fields{
			"cacheMB":    cacheMB,
			"cacheTTL":   cacheTTL,
			"ListenPort": GinPort,
		}).Print("Gin Starting now")

		if err := runGin(
			GinPort,
			mtcache.Config{
				MaxSizeMB: int64(cacheMB),
				MaxTTL:    cacheTTL,
			},
			map[mtproxy.RoutingKey]string{
				{SourceLang: "en", TargetLang: "pt"}: "bananas.foo",
				{}:                                   "bar.foo",
			},
		); err != nil {
			log.Error("ServeHTTP error", err)
		}

	},
}

func runGin(
	listenPort int16, cacheCfg mtcache.Config, routes map[mtproxy.RoutingKey]string,
) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		ginlogrus.Logger(log),
		gin.Recovery(),
	)
	r.Use(gzip.Gzip(gzip.BestSpeed, gzip.WithExcludedPaths([]string{"/metrics"})))

	// TODO: cmdline args for this
	srv, err := server.NewGinServer(
		cacheCfg,
		routes,
	)
	if err != nil {
		return err
	}

	r.GET("/ping", srv.Ping)
	r.GET("/echo/:id/:cnt", srv.Message)
	r.POST("/v1/machine_translate", srv.MachineTranslate)

	return r.Run(fmt.Sprintf(":%v", listenPort))
}
