package handler

import (
	"fmt"

	"github.com/msf/cachingproxy/model"
	log "github.com/sirupsen/logrus"
)

type cachingMTHandler struct {
	localCache       SegmentRepositoryCache
	remoteTranslator SegmentRepository
}

type SegmentRepositoryCache interface {
	SegmentRepository
	Save(model.MTRequestMetadata, []model.TargetSegment)
	Metrics() string
}

type SegmentRepository interface {
	Handle(*model.MachineTranslationRequest) ([]model.TargetSegment, error)
}

func NewCachingMTHandler() *cachingMTHandler {
	return &cachingMTHandler{}

}

func (m *cachingMTHandler) Handle(req *model.MachineTranslationRequest) (resp *model.MachineTranslationResponse, err error) {
	if err := req.HasError(); err != nil {
		return resp, fmt.Errorf("handler: invalid message: %#v, %v", req, err)
	}

	resp = &model.MachineTranslationResponse{
		RequestID:       req.ID,
		RequestMetadata: req.Metadata,
		TargetSegments:  make([]model.TargetSegment, len(req.Segments)),
	}
	resp.RequestID = req.ID

	log.WithFields(log.Fields{
		"id":           req.ID,
		"sourceLang":   req.Metadata.SourceLang,
		"targetLang":   req.Metadata.TargetLang,
		"segmentCount": len(req.Segments),
	}).Info("MachineTranslate")

	// fetch from cache
	tSegments, err := m.localCache.Handle(req)

	// find what we're missing
	missingSources := make([]string, len(tSegments)/2)
	missingIndexes := make([]int, len(tSegments)/2)
	for i, v := range tSegments {
		if v == "" && req.Segments[i] != "" {
			missingIndexes = append(missingIndexes, i)
			missingSources = append(missingSources, req.Segments[i])
		} else {
			resp.TargetSegments[i] = v
		}
	}

	// get the missing segments
	rSegments, err := m.remoteTranslator.Handle(&model.MachineTranslationRequest{
		ID:       req.ID,
		Metadata: req.Metadata,
		Segments: missingSources,
	})

	m.localCache.Save(req.Metadata, rSegments)

	// finish the response results
	for i, v := range rSegments {
		pos := missingIndexes[i]
		resp.TargetSegments[pos] = v
	}

	//TODO: emit metrics to prometheus
	log.WithFields(log.Fields{
		"hitCount":  len(tSegments) - len(missingIndexes),
		"missCount": len(missingIndexes),
		"metrics":   m.localCache.Metrics(),
	}).Info("Translation Complete")

	return resp, nil
}
