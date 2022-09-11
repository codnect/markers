package processor

import "sync"

type CommandCallback func(ctx *Context)

var (
	callbackMu       sync.Mutex
	generateCallback CommandCallback
	validateCallback CommandCallback
)

func SetGenerateCommandCallback(callback CommandCallback) {
	defer callbackMu.Unlock()
	callbackMu.Lock()
	generateCallback = callback
}

func getGenerateCommandCallback() CommandCallback {
	return generateCallback
}

func SetValidateCommandCallback(callback CommandCallback) {
	defer callbackMu.Unlock()
	callbackMu.Lock()
	validateCallback = callback
}

func getValidateCommandCallback() CommandCallback {
	return validateCallback
}
