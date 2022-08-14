package processor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"sync"
)

type GenerateCallback func(collector *marker.Collector, loadResult *packages.LoadResult, dirs []string) error

var (
	generateCallbackMu sync.RWMutex
	generateCallback   GenerateCallback
)

func setGenerateCallback(callback GenerateCallback) {
	generateCallbackMu.Lock()
	defer generateCallbackMu.Unlock()
	generateCallback = callback
}

func invokeGenerateCallback(collector *marker.Collector, loadResult *packages.LoadResult, dirs []string) error {
	return generateCallback(collector, loadResult, dirs)
}
