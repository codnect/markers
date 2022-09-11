package main

import (
	"github.com/procyon-projects/marker/processor"
	"log"
)

func init() {
	processor.Initialize(Package, AppName, Version)
	processor.SetGenerateCommandCallback(Generate)
	processor.SetValidateCommandCallback(Validate)
	processor.AddRegistryFunction(nil)
}

func main() {
	log.SetFlags(0)
	processor.Execute()
}
