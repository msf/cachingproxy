package handler

import (
	"fmt"

	"github.com/msf/cachingproxy/model"
	log "github.com/sirupsen/logrus"
)

type cachingMTHandler struct {
	localCache       SegmentRepository
	remoteTranslator SegmentRepository
}

type SegmentRepository interface {
	Translate([]string) (map[int]string, error)
}

func New() *cachingMTHandler {

}

func MachineTranslate(req model.MachineTranslationRequest) (resp model.MachineTranslationResponse, err error) {
	if err := req.HasError(); err != nil {
		return resp, fmt.Errorf("handler: invalid message: %#v, %v", req, err)
	}

	log.WithFields(log.Fields{
		"id":           req.ID,
		"segmentCount": len(req.Segments),
	}).Info("MachineTranslate")

	return resp, nil
}
