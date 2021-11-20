// Package server provides an application server
package server

import (
	"context"
	"fmt"

	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
	api "github.com/msf/cachingproxy/proto/gen/go/application/api/v1"
)

// EchoServiceServer type
type EchoServiceServer struct {
}

// Echo echoes
func (e *EchoServiceServer) Echo(ctx context.Context, r *api.EchoRequest) (resp *api.EchoResponse, err error) {

	m := model.Message{ID: r.Id, Content: r.Content}
	if !m.Valid() {
		return nil, fmt.Errorf("invalid request id: %s, content: %s", m.ID, m.Content)
	}

	if m, err = handler.EchoMessage(m); err != nil {
		return nil, err
	}

	resp = &api.EchoResponse{Id: m.ID, Content: m.Content}
	return resp, nil
}

// EchoStream continues
func (e *EchoServiceServer) EchoStream(r *api.EchoStreamRequest, ss api.EchoService_EchoStreamServer) (err error) {

	for x := 0; x < len(r.Content); x++ {
		resp, err := e.Echo(ss.Context(), &api.EchoRequest{Id: r.Id, Content: r.Content})
		if err != nil {
			return err
		}

		if err := ss.Send(&api.EchoStreamResponse{Id: resp.Id, Content: resp.Content}); err != nil {
			return err
		}

	}

	return nil
}
