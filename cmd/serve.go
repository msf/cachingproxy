package cmd

import (
	"context"
	"os"

	"github.com/msf/cachingproxy/proto/gen/go/application/api/v1"
	"github.com/msf/cachingproxy/server"
	"github.com/spf13/cobra"
	"gitlab.com/brunotm/monorepo/pkg/cobrautil"
	"gitlab.com/brunotm/monorepo/pkg/gserver"
	"gitlab.com/brunotm/monorepo/pkg/gserver/interceptor"
	"gitlab.com/brunotm/monorepo/pkg/log"
	"gitlab.com/brunotm/monorepo/pkg/sigaction"
	"go.uber.org/zap"
)

// Serve starts the application server
var Serve = &cobra.Command{
	Use:   "serve",
	Short: "start serving",
	Run: func(cmd *cobra.Command, args []string) {
		log.ReplaceGrpcLogger()

		c := &gserver.Config{}
		c.Gateway.Enabled = true
		c.Gateway.InProcess = true
		c.GrpcAddr = "127.0.0.1:9000"
		c.HTTPAddr = "127.0.0.1:8080"
		c.Interceptors = &interceptor.Config{}
		c.Interceptors.Recovery = true
		c.Interceptors.ContextTags.Enabled = true
		c.Interceptors.Metrics.Enabled = true
		c.Interceptors.Metrics.HistogramBuckets = []float64{0.5, 0.75, 0.90, 0.99}
		c.Interceptors.Metrics.Path = "/metrics"

		// Enable development features.
		if cobrautil.GetFlagBool(cmd, "dev") {
			c.GoDebug = true
			c.Reflection = true
		}

		srv, err := gserver.New(c, log.Root())
		if err != nil {
			log.Root().Fatal("error creating server", zap.Error(err))
		}

		// Register the application service with the server
		srv.RegisterService(&api.EchoService_ServiceDesc, &server.EchoServiceServer{})

		// Register our client with the grpc-gateway server
		if err = srv.RegisterGatewayClient(api.RegisterEchoServiceHandler); err != nil {
			log.Root().Fatal("error registering gateway client", zap.Error(err))
		}

		// Gracefully stop on SIGINT
		go sigaction.Handle(func(sig os.Signal) bool {
			log.Root().Info("stopping server", zap.String("signal", sig.String()))
			if err := srv.Close(context.Background()); err != nil {
				log.Root().Error("error stopping server", zap.Error(err))
			}

			os.Exit(0)
			return false
		}, sigaction.SIGINT)

		// Start the server
		log.Root().Info("starting server")
		if err = srv.ListenAndServe(); err != nil {
			log.Root().Fatal("error running server", zap.Error(err))
		}
	},
}
