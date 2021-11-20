package handler

import (
	"fmt"
	"strings"

	"github.com/msf/cachingproxy/model"
	log "github.com/sirupsen/logrus"
)

// EchoMessage generates a new message with the same content
func EchoMessage(m model.Message) (e model.Message, err error) {
	if !m.Valid() {
		return e, fmt.Errorf("handler: invalid message: %#v", m)
	}

	payload := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		payload = append(payload, m.Content)
	}
	e.ID = "42"
	e.Content = strings.Join(payload, "-")

	log.WithFields(log.Fields{
		"id":      m.ID,
		"content": m.Content,
	}).Info("Echo")

	return e, nil
}
