package logpeck

import (
	"log"
	"os"
	"time"
)

func LogExecTime(start time.Time, prefix string) {
	elapsed_ms := time.Since(start) / time.Millisecond
	log.Printf("Performance: %s cost %d ms.\n", prefix, elapsed_ms)
}

func GetHost() string {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return host
}
