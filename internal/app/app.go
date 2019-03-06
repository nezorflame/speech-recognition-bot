package app

import (
	"log"
)

// Cmd vars, overridden by ldflags
var (
	Name      = "_"
	Version   = "devel"
	BuildTS   = "_"
	GoVersion = "_"
	GitHash   = "_"
	GitBranch = "_"
)

// PrintInfo logs the inforamtion about the current launched binary if debug mode is enabled
func PrintInfo(debug bool) {
	if !debug {
		return
	}

	log.Println("App info:")
	log.Printf("  Version: %s", Version)
	log.Printf("  Build date: %s", BuildTS)
	log.Printf("  Go version: v%s", GoVersion)

	log.Println("Git info:")
	log.Printf("  Tag: %s", GitBranch)
	log.Printf("  Commit: %s", GitHash)
}
