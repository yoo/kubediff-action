package main

import (
	"log"
	"os"
)

var debug bool

func init() {
	if os.Getenv("INPUT_DEBUG") != "" {
		debug = true
	}
}

func Debugf(format string, v ...interface{}) {
	if !debug {
		return
	}
	log.Printf("DEBUG: "+format, v...)
}
