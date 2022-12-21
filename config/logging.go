package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/syhily/hobbit/pkg/toml"
)

var (
	// defaultParentDir is the default directory of hobbit.
	defaultParentDir = filepath.Join(".", "data")
)

const loggingTmpl = `
## logging related configuration.
[logging]

## Dir is the output directory for log-files
## Default: %s
dir = "%s"

## Determine which level of logs will be emitted.
## error, warn, info, and debug are available
## Default: %s
level = "%s"

## MaxSize is the maximum size in megabytes of the log file before it gets rotated.
## Default: %s
maxSize = "%s"

## MaxBackups is the maximum number of old log files to retain. The default
## is to retain all old log files (though MaxAge may still cause them to get deleted.)
## Default: %d
maxBackups = %d

## MaxAge is the maximum number of days to retain old log files based on the
## timestamp encoded in their filename.  Note that a day is defined as 24 hours
## and may not exactly correspond to calendar days due to daylight savings, leap seconds, etc.
## The default is not to remove old log files based on age.
## Default: %d
maxAge = %d`

// Logging represents a logging configuration
type Logging struct {
	Dir        string    `toml:"dir"`
	Level      string    `toml:"level"`
	MaxSize    toml.Size `toml:"maxSize"`
	MaxBackups uint16    `toml:"maxBackups"`
	MaxAge     uint16    `toml:"maxAge"`
}

func (l *Logging) TOML() string {
	dl := NewDefaultLogging()

	return fmt.Sprintf(
		loggingTmpl,
		dl.Dir,
		strings.ReplaceAll(l.Dir, `\`, `\\`),
		dl.Level,
		l.Level,
		&dl.MaxSize,
		&l.MaxSize,
		dl.MaxBackups,
		l.MaxBackups,
		dl.MaxAge,
		l.MaxAge,
	)
}

// NewDefaultLogging returns a new default logging config
func NewDefaultLogging() *Logging {
	return &Logging{
		Dir:        filepath.Join(defaultParentDir, "log"),
		Level:      "info",
		MaxSize:    toml.Size(100 * 1024 * 1024),
		MaxBackups: 3,
		MaxAge:     7,
	}
}
