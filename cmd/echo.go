package cmd

import (
	"fmt"

	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
	"github.com/spf13/cobra"
	"gitlab.com/brunotm/monorepo/pkg/cobrautil"
	"gitlab.com/brunotm/monorepo/pkg/log"
	"go.uber.org/zap"
)

// Run uses the application handler
var Run = &cobra.Command{
	Use:   "echo",
	Short: "echoes any given message input",
	Run: func(cmd *cobra.Command, args []string) {

		id := cobrautil.GetFlagString(cmd, "id")
		content := cobrautil.GetFlagString(cmd, "content")
		m, err := handler.EchoMessage(model.Message{ID: id, Content: content})

		if err != nil {
			log.Root().Fatal("error running echo", zap.Error(err))
		}

		fmt.Printf("echo response id: %s, content: %s\n", m.ID, m.Content)
	},
}
