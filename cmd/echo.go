package cmd

import (
	"compress/gzip"
	"fmt"

	"strings"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/msf/cachingproxy/handler/mtcache"
	"github.com/msf/cachingproxy/handler/mtproxy"
	"github.com/msf/cachingproxy/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	EchoPort int16
)

func init() {
	rootCmd.AddCommand(echoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// echoCmd.PersistentFlags().String("foo", "", "A help for foo")

	ginCmd.PersistentFlags().Int16Var(&EchoPort, "echoPort", 4322, "gin server listening port")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	echoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// echoCmd represents the echo command
var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "mtproxy using echo",
	Run: func(cmd *cobra.Command, args []string) {

		log = logrus.New()
		log.WithFields(logrus.Fields{
			"cacheMB":    cacheMB,
			"cacheTTL":   cacheTTL,
			"ListenPort": EchoPort,
		}).Print("Echo Starting now")

		if err := runEcho(
			EchoPort,
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

func runEcho(
	listenPort int16, cacheCfg mtcache.Config, routes map[mtproxy.RoutingKey]string,
) error {
	e := echo.New()
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzip.BestSpeed,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/metrics")
		},
	}))

	log.WithFields(logrus.Fields{
		"cacheConfig": cacheCfg,
		"routes":      routes,
	}).Info(("Echo is ready!"))

	e.GET("/echo/:id/:cnt", server.EchoMessage)
	e.GET("/ping", server.EchoPing)
	return e.Start(fmt.Sprintf(":%v", listenPort))
}
