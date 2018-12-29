package logpeck

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// LogExecTime .
func LogExecTime(start time.Time, prefix string) {
	elapsedMs := time.Since(start) / time.Millisecond
	log.Debugf("Performance: %s cost %d ms.", prefix, elapsedMs)
}

// GetHost .
func GetHost() string {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return host
}

// SelectRandom .
func SelectRandom(candidates []string) (string, error) {
	l := len(candidates)
	if l <= 0 {
		return "", errors.New("none candidates")
	}
	ret := candidates[rand.Intn(l)]
	return ret, nil
}

// SplitString .
func SplitString(content, delims string) []string {
	if len(delims) == 0 {
		delims = "\t\r\n "
	}
	splitFunc := func(r rune) bool {
		for _, d := range delims {
			if r == d {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(content, splitFunc)
}
