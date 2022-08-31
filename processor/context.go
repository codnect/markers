package processor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
)

type Context struct {
	collector      *marker.Collector
	loadResult     *packages.LoadResult
	config         Config
	configFilePath string
}

func (ctx *Context) Config() Config {
	return ctx.config
}

func (ctx *Context) ConfigFilePath() string {
	return ctx.configFilePath
}

func (ctx *Context) Collector() *marker.Collector {
	return ctx.collector
}

func (ctx *Context) LoadResult() *packages.LoadResult {
	return ctx.loadResult
}
