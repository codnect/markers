package processor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"log"
)

type Context struct {
	dirs           []string
	loadResult     *packages.LoadResult
	registry       *marker.Registry
	collector      *marker.Collector
	config         Config
	configFilePath string
}

func (ctx *Context) Directories() []string {
	return ctx.dirs
}

func (ctx *Context) Registry() *marker.Registry {
	return ctx.registry
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

func (ctx *Context) PrintError(err error) {
	if err != nil {

		switch typedErr := err.(type) {
		case marker.ErrorList:
			ctx.PrintErrors(typedErr)
			return
		}

		log.Println(err)
		return
	}
}

func (ctx *Context) PrintErrors(errorList marker.ErrorList) {
	if errorList == nil || len(errorList) == 0 {
		return
	}

	for _, err := range errorList {
		switch typedErr := err.(type) {
		case marker.Error:
			pos := typedErr.Position
			log.Printf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ParserError:
			pos := typedErr.Position
			log.Printf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ErrorList:
			ctx.PrintErrors(typedErr)
		default:
			ctx.PrintError(err)
		}
	}
}
