package config

import (
	"fmt"
	"runtime"
)

// These variables are populated via the Go linker.
var (
	GitVersion = "" // Semantic version, composed of ${Major}.${Minor}.${Patch}
	GitCommit  = "" // sha1 from git, output of $(git rev-parse HEAD)
	BuildTime  = "" // build time in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	GoVersion  = runtime.Version()
	Platform   = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)
