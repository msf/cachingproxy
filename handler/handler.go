// Package handler implements application main business logic
package handler

import "github.com/msf/cachingproxy/model"

type MachineTranslationHandler interface {
	Handle(*model.MachineTranslationRequest) (*model.MachineTranslationResponse, error)
}
