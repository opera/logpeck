package logpeck

import (
	"log"
	"time"
)

func LogExecTime(start time.Time, prefix string) {
	elapsed_ms := time.Since(start) / time.Millisecond
	log.Printf("Performance: %s cost %d ms.\n", prefix, elapsed_ms)
}
