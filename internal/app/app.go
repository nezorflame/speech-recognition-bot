package app

import (
	log "github.com/sirupsen/logrus"
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

// PrintInfo logs the inforamtion about the current launched binary if debug mode is enabled)
func PrintInfo() {
	log.Debug("App info:")
	log.Debugf("  Version: %s", Version)
	log.Debugf("  Build date: %s", BuildTS)
	log.Debugf("  Go version: v%s", GoVersion)

	log.Debug("Git info:")
	log.Debugf("  Tag: %s", GitBranch)
	log.Debugf("  Commit: %s", GitHash)
}
