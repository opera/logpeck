package logpeck

import (
	"log"
	"testing"
)

func TestGetHost(t *testing.T) {
	log.Println("local host: " + GetHost())
}
