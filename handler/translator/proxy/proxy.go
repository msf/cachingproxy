package proxy

import (
	"net/http"

	"github.com/msf/cachingproxy/model"
)

type MaestroSegmentTranslator struct {
	client http.Client

	// used to identify to which hostname/path a request should go
	routingMap map[string]string
}

func (m *MaestroSegmentTranslator) Handle(req model.MachineTranslationRequest) (resp model.MachineTranslationResponse, err error) {
}
