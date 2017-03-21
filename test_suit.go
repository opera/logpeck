package logpeck

import (
	"fmt"
	"time"
)

const kTestDBPath string = ".unittest.db"

func logExecTime(start time.Time, prefix string) {
	elapsed_ms := time.Since(start) / time.Millisecond
	fmt.Printf("Performance: %s cost %d ms.\n", prefix, elapsed_ms)
}
