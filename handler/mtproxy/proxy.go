package mtproxy

import (
	"fmt"
	"net/http"

	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
)

type RoutingKey struct { // TODO use more fields
	SourceLang string
	TargetLang string
}

// MaestroProxyTranslator translates by calling maestro endpoints
type MaestroProxyTranslator struct {
	client http.Client // TODO http retry logic

	// used to identify to which hostname/path a request should go
	routingMap map[RoutingKey]string
}

func NewMaestroProxyTranslator(routingMap map[RoutingKey]string) (handler.MachineTranslationHandler, error) {
	if len(routingMap) < 1 {
		return nil, fmt.Errorf("MaestroProxyTranslator needs a routingMap, got zero entries")
	}
	return &MaestroProxyTranslator{
		client:     *http.DefaultClient,
		routingMap: routingMap,
	}, nil
}

func (m *MaestroProxyTranslator) Handle(
	req *model.MachineTranslationRequest) (resp *model.MachineTranslationResponse, err error) {
	k := keyForReq(req)
	hostname, found := m.routingMap[k]
	if !found {
		err = fmt.Errorf("no hostname for %+v for MaestroProxyTranslator", k)
		return
	}
	resp, err = m.doRequest(hostname, req)
	return
}

func keyForReq(req *model.MachineTranslationRequest) RoutingKey {
	return RoutingKey{
		SourceLang: req.Metadata.SourceLang,
		TargetLang: req.Metadata.TargetLang,
	}
}

func (m *MaestroProxyTranslator) doRequest( //TODO: invoke maestro
	hostname string, req *model.MachineTranslationRequest,
) (resq *model.MachineTranslationResponse, err error) {
	err = fmt.Errorf("doRequest not implemented")
	return
}
