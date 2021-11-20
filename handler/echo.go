package handler

import (
	"fmt"

	"github.com/msf/cachingproxy/model"
	"gitlab.com/brunotm/monorepo/pkg/log"
	"gitlab.com/brunotm/monorepo/pkg/random"
	"go.uber.org/zap"
)

const idLength = 10

// EchoMessage generates a new message with the same content
func EchoMessage(m model.Message) (e model.Message, err error) {
	if !m.Valid() {
		return e, fmt.Errorf("handler: invalid message: %#v", m)
	}

	e.ID = random.String(idLength)
	e.Content = m.Content

	log.Root().Info("handled message", zap.String("id", m.ID))
	return e, nil
}
