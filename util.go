package logpeck

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"os"
	"strings"
	"time"
)

func LogExecTime(start time.Time, prefix string) {
	elapsed_ms := time.Since(start) / time.Millisecond
	log.Debugf("Performance: %s cost %d ms.\n", prefix, elapsed_ms)
}

func GetHost() string {
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return host
}

func SelectRandom(candidates []string) (string, error) {
	candi_len := len(candidates)
	if candi_len <= 0 {
		return "", errors.New("none candidates")
	}
	ret := candidates[rand.Intn(candi_len)]
	return ret, nil
}

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
