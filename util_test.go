package logpeck

import (
	"log"
	"testing"
)

func TestGetHost(t *testing.T) {
	log.Println("local host: " + GetHost())
}

func TestSplitString(t *testing.T) {
	content := "hello world, golang"
	delims := " ,"
	arr := SplitString(content, delims)
	if len(arr) != 3 {
		panic(arr)
	}
	log.Println(arr)

	content = `
		"Name":"TestLog",
		"LogPath":".test.log",
		"Fields":[
		{
			"Name": "Date",
			"Type": "string",
			"Value": "$1"
		}
	}`
	delims = "\n\r\t\" ,:{}[]"
	arr = SplitString(content, delims)
	log.Println(arr)
	if len(arr) != 11 {
		panic(len(arr))
	}
}
