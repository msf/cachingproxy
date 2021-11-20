package handler

import (
	"fmt"

	"github.com/msf/cachingproxy/model"
	log "github.com/sirupsen/logrus"
)

// EchoMessage generates a new message with the same content
func EchoMessage(m model.Message) (e model.Message, err error) {
	if !m.Valid() {
		return e, fmt.Errorf("handler: invalid message: %#v", m)
	}

	e.ID = "42"
	e.Content = m.Content

	log.WithFields(log.Fields{
		"content": m.Content,
	}).Info("Echo done")
	return e, nil
}
