package mtcache

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/msf/cachingproxy/handler"
	"github.com/msf/cachingproxy/model"
)

type Config struct {
	MaxSizeMB int64
	MaxTTL    time.Duration
}

type MachineTranslationCache interface {
	handler.MachineTranslationHandler
	Save(model.MTRequestMetadata, []string, []model.TargetSegment) error
	Metrics() string
}

type machineTranslationCache struct {
	cache  *ristretto.Cache
	config Config
}

func NewCachingSegmentTranslator(config Config) (MachineTranslationCache, error) {
	// docs of ristretto just say use 64  :-|
	const BufferItemCount = 64
	cache, err := ristretto.NewCache(&ristretto.Config{
		MaxCost:     config.MaxSizeMB << 20, // mb to bytes
		BufferItems: BufferItemCount,
		// assuming ~100bytes per entry
		NumCounters: 10_000 * config.MaxSizeMB,
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}
	return &machineTranslationCache{
		cache:  cache,
		config: config,
	}, nil
}

func (c *machineTranslationCache) Handle(
	req *model.MachineTranslationRequest,
) (*model.MachineTranslationResponse, error) {
	keys := keysFor(req.Metadata, req.Segments)

	resp := make([]model.TargetSegment, len(keys))
	for i, k := range keys {
		var val model.TargetSegment
		v, found := c.cache.Get(k)
		if !found {
			// empty strings to indicate cache miss
			val = ""
		} else {
			val = v.(model.TargetSegment)
		}
		resp[i] = val
	}
	return &model.MachineTranslationResponse{
		RequestID:       req.ID,
		TargetSegments:  resp,
		RequestMetadata: req.Metadata,
	}, nil
}

func (c *machineTranslationCache) Save(
	metadata model.MTRequestMetadata,
	sourceSegments []string,
	targetSegments []model.TargetSegment,
) error {

	if len(sourceSegments) != len(targetSegments) {
		return fmt.Errorf("non matching source and target segment array lengths")
	}

	keys := keysFor(metadata, sourceSegments)
	for i, seg := range targetSegments {
		c.cache.SetWithTTL(
			keys[i],
			seg,
			int64(len([]byte(seg))),
			c.config.MaxTTL,
		)
	}
	return nil
}

func (c *machineTranslationCache) Metrics() string {
	return c.cache.Metrics.String()
}

func keysFor(md model.MTRequestMetadata, segments []string) []string {
	var b strings.Builder
	b.WriteString(md.SourceLang)
	b.WriteString("|")
	b.WriteString(md.TargetLang)
	b.WriteString("|")
	//TODO: fmt is non-ideal on hot-paths
	b.WriteString(fmt.Sprintf("%+v", md.Metadata))
	b.WriteString("|")
	prefix := b.String()

	keys := make([]string, len(segments))
	for _, v := range segments {
		keys = append(keys, prefix+v)
	}
	return keys
}
