// Package main implements the application initialization code
package main

import (
	"github.com/msf/cachingproxy/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	// Name injected at build time
	Name string

	// Version injected at build time
	Version string

	// BuildTime injected at build time
	BuildTime string
)

// init is only used for keeping the command setup within the same file
func main() {
	log.WithFields(log.Fields{
		"Name":      Name,
		"Version":   Version,
		"BuildTime": BuildTime,
	}).Print("Starting now")

	if err := cmd.ServeHTTP(); err != nil {
		log.Error("ServeHTTP error", err)
	}
}
