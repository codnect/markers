package main

import (
	"github.com/procyon-projects/marker/processor"
	"log"
)

func init() {
	processor.Initialize(AppName, AppVersion)
}

func main() {
	log.SetFlags(0)
	processor.Execute()
}
