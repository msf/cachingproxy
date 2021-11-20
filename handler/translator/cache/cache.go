package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/msf/cachingproxy/model"
)

type Config struct {
	MaxSizeKB int
	MaxTTL    time.Duration
}

type CachingSegmentTranslator struct {
	cache *ristretto.Cache
}

func (c *CachingSegmentTranslator) Save(req model.MachineTranslationRequest, resp model.MachineTranslationResponse) {
}

func (m *CachingSegmentTranslator) Handle(req model.MachineTranslationRequest) 
(resp model.MachineTranslationResponse, err error) {
	keys := keysFor(req.Metadata, req.Segments)

	missing := make([]int, len(keys)
	found := make(map[int]string, len(keys))
	for i, k range keys {
		val, found := m.cache.Get(k)
		if !found {
			missing = append(missing, i)
		} else {
			found[i] = String(val)
		}
	}
	val, found := m.cache.Get()
}

func keyFor(md model.MTRequestMetadata, segments []string) map[string]string {
	return map[string]string{}
}
