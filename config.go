package logpeck

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// LPConfig .
type LPConfig struct {
	Port          int32         `toml:"port"`
	LogLevel      string        `toml:"log_level"`
	MaxTaskNum    int32         `toml:"max_task_num"`
	DatabaseFile  string        `toml:"database_file"`
	PeckTaskLimit PeckTaskLimit `toml:"peck_task_limit"`
}

// PeckTaskLimit .
type PeckTaskLimit struct {
	MaxLinesPerSec int64 `toml:"max_lines_per_sec"`
	MaxBytesPerSec int64 `toml:"max_bytes_per_sec"`
}

// Config global variable
var Config LPConfig

// InitConfig init Config
func InitConfig(file *string) bool {
	Config = LPConfig{
		Port:         7117,
		MaxTaskNum:   16,
		DatabaseFile: "logpeck.db",
	}

	if _, err := toml.DecodeFile(*file, &Config); err != nil {
		fmt.Fprintf(os.Stderr, "Parse config fail: %s.\n", err)
		return false
	}
	return true
}
