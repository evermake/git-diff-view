package v1

import (
	"github.com/evermake/git-diff-view/internal/controller/http/v1/openapi"
	"github.com/evermake/git-diff-view/pkg/diff"
)

type diffCache struct {
	internal map[string][]*diffCacheEntry
}

func newDiffCache() *diffCache {
	return &diffCache{
		internal: make(map[string][]*diffCacheEntry),
	}
}

func (c *diffCache) hash(commitA, commitB string) string {
	return commitA + commitB
}

func (c *diffCache) Get(commitA, commitB string) ([]*diffCacheEntry, bool) {
	e, ok := c.internal[c.hash(commitA, commitB)]
	return e, ok
}

func (c *diffCache) Set(commitA, commitB string, entry []*diffCacheEntry) {
	c.internal[c.hash(commitA, commitB)] = entry
}

type diffCacheEntry struct {
	FileDiff *openapi.FileDiff
	Diff     *diff.Diff
}