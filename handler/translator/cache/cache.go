package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/msf/cachingproxy/model"
)

type Config struct {
	MaxSizeMB int64
	MaxTTL    time.Duration
}

type CachingSegmentTranslator struct {
	cache  *ristretto.Cache
	config Config
}

func NewCachingSegmentTranslator(config Config) (*CachingSegmentTranslator, error) {
	// docs of ristretto just say use 64  :-|
	const BufferItemCount = 64
	cache, err := ristretto.NewCache(&ristretto.Config{
		MaxCost:     int64(config.MaxSizeMB << 20),
		BufferItems: BufferItemCount,
		// assuming ~100bytes per entry
		NumCounters: 10_000 * config.MaxSizeMB,
		Metrics:     true,
	})
	if err != nil {
		return nil, err
	}
	return &CachingSegmentTranslator{
		cache:  cache,
		config: config,
	}, nil
}

func (c *CachingSegmentTranslator) Save(
	metadata model.MTRequestMetadata,
	sourceSegments []string,
	targetSegments []model.TargetSegment,
) error {

	if len(sourceSegments) != len(targetSegments) {
		return fmt.Errorf("non matching source and target segment array lengths")
	}

	keys := keysFor(metadata, sourceSegments)
	for i, tgt_seg := range targetSegments {
		c.cache.SetWithTTL(
			keys[i],
			tgt_seg,
			int64(len([]byte(tgt_seg))),
			c.config.MaxTTL,
		)
	}
	return nil
}

func (m *CachingSegmentTranslator) Handle(
	req *model.MachineTranslationRequest,
) ([]model.TargetSegment, error) {
	keys := keysFor(req.Metadata, req.Segments)

	resp := make([]model.TargetSegment, len(keys))
	for i, k := range keys {
		var val model.TargetSegment
		v, found := m.cache.Get(k)
		if !found {
			// empty strings indicate cache miss
			val = ""
		} else {
			val = v.(model.TargetSegment)
		}
		resp[i] = val
	}
	return resp, nil
}

func (m *CachingSegmentTranslator) Metrics() string {
	return m.cache.Metrics.String()
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
