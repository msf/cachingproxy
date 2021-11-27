package handler

import (
	"fmt"

	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/handler/mtcache"
	"github.com/msf/cachingproxy/handler/mtproxy"
	"github.com/msf/cachingproxy/model"
	log "github.com/sirupsen/logrus"
)

type cachingMTHandler struct {
	localCache       mtcache.MachineTranslationCache
	remoteTranslator handler.MachineTranslationHandler
}

func NewCachingMTHandler(
	cacheConfig mtcache.Config, routingMap map[mtproxy.RoutingKey]string,
) (handler.MachineTranslationHandler, error) {

	cache, err := mtcache.NewCachingSegmentTranslator(cacheConfig)
	if err != nil {
		return nil, err
	}
	remote, err := mtproxy.NewMaestroProxyTranslator(routingMap)
	if err != nil {
		return nil, err
	}
	return &cachingMTHandler{
		localCache:       cache,
		remoteTranslator: remote,
	}, nil
}

func (m *cachingMTHandler) Handle(
	req *model.MachineTranslationRequest,
) (resp *model.MachineTranslationResponse, err error) {
	if err := req.HasError(); err != nil {
		return resp, fmt.Errorf("handler: invalid message: %#v, %v", req, err)
	}

	log.WithFields(log.Fields{
		"id":           req.ID,
		"sourceLang":   req.Metadata.SourceLang,
		"targetLang":   req.Metadata.TargetLang,
		"segmentCount": len(req.Segments),
	}).Info("MachineTranslate")

	// fetch from cache
	resp, err = m.localCache.Handle(req)
	if err != nil {
		log.Error("cache req failed", err)
		// TODO more metrics
		return resp, err
	}

	// find what we're missing
	hitCount := len(req.Segments)
	missingSources := make([]string, len(resp.TargetSegments)/2)
	missingIndexes := make([]int, len(resp.TargetSegments)/2)
	for i, v := range resp.TargetSegments {
		if v == "" && req.Segments[i] != "" {
			hitCount--
			missingIndexes = append(missingIndexes, i)
			missingSources = append(missingSources, req.Segments[i])
		}
	}

	// get the missing segments
	rResp, err := m.remoteTranslator.Handle(&model.MachineTranslationRequest{
		ID:       req.ID,
		Metadata: req.Metadata,
		Segments: missingSources,
	})
	if err != nil {
		log.Error("remoteTranslator failed", err)
		// TODO more metrics
		return resp, err
	}

	er := m.localCache.Save(req.Metadata, req.Segments, rResp.TargetSegments)
	if er != nil {
		log.Error("locaCache.Save() failed", er)
	}

	// finish the response results
	for i, v := range rResp.TargetSegments {
		pos := missingIndexes[i]
		resp.TargetSegments[pos] = v
	}

	//TODO: emit metrics to prometheus
	log.WithFields(log.Fields{
		"hitCount":  hitCount,
		"missCount": len(missingIndexes),
		"metrics":   m.localCache.Metrics(),
	}).Info("Translation Complete")

	return resp, nil
}
