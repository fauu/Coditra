package coditra

import "github.com/twmb/murmur3"

type LookupRequestHash = uint64

type LookupCache struct {
	entries map[LookupRequestHash]any
}

func NewLookupCache() LookupCache {
	return LookupCache{
		entries: make(map[LookupRequestHash]any),
	}
}

// QUAL: Make the input params into a `LookupRequest` struct?
func (c LookupCache) query(sourceID string, input string) (LookupRequestHash, *any) {
	requestHash := hashRequest(sourceID, input)
	result, found := c.entries[requestHash]
	if found {
		return requestHash, &result
	}
	return requestHash, nil
}

func (c LookupCache) store(requestHash LookupRequestHash, result any) {
	c.entries[requestHash] = result
}

func hashRequest(sourceID string, input string) LookupRequestHash {
	var hash = murmur3.New64()
	hash.Write([]byte(sourceID + input))
	return hash.Sum64()
}
