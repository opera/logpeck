package logpeck

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type PeckerConfig struct {
	Port       int32     `toml:"port"`
	MaxTaskNum int32     `toml:"max_task_num"`
	TaskLimit  TaskLimit `toml:"task_limit"`
}

type TaskLimit struct {
	MaxLinesPerSec int64 `toml:"max_lines_per_sec"`
	MaxBytesPerSec int64 `toml:"max_bytes_per_sec"`
}

var PkConfig PeckerConfig

func InitConfig(file string) bool {
	if _, err := toml.DecodeFile(file, &PkConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Parse config fail: %s.\n", err)
		return false
	}
	return true
}
