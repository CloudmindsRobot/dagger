package databases

import "github.com/bluele/gcache"

var (
	GC gcache.Cache
)

func init() {
	GC = gcache.New(20).LRU().Build()
}
