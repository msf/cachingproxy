// Package main implements the application initialization code
package main

import (
	"fmt"

	"github.com/msf/cachingproxy/cmd"
	"github.com/spf13/cobra"
	"gitlab.com/brunotm/monorepo/pkg/cobrautil"
	"gitlab.com/brunotm/monorepo/pkg/log"
	"go.uber.org/zap"
)

var (
	// Name injected at build time
	Name string

	// Version injected at build time
	Version string

	// BuildTime injected at build time
	BuildTime string
)

// command is the main application command
var command = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		_ = log.SetLevel(cobrautil.GetFlagString(cmd, "log-level"))
	},

	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			log.Root().Fatal("error running command", zap.Error(err))
		}
	},
}

// init is only used for keeping the command setup within the same file
func main() {
	command.Use = Name
	command.Version = fmt.Sprintf("%s, build time %s", Version, BuildTime)

	command.PersistentFlags().String("log-level", "INFO", "log level")

	command.AddCommand(cmd.Run)
	cmd.Run.Flags().String("id", "", "message id for echoing")
	cmd.Run.Flags().String("content", "", "message content for echoing")

	command.AddCommand(cmd.Serve)
	cmd.Serve.Flags().Bool("gateway", true, "serve http")

	if err := command.Execute(); err != nil {
		log.Root().Fatal("error running command", zap.Error(err))
	}
}
