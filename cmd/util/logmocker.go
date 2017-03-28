package main

import (
	"flag"
	"fmt"
	"github.com/opera/logpeck"
)

func Usage() {
	s := "Usage: \n\t./logmocker --filename=<logfile>"
	fmt.Println(s)
}
func main() {
	filename := flag.String("filename", "", "log file name")
	flag.Parse()
	if *filename == "" {
		Usage()
		return
	}
	mocker, err := logpeck.NewMockLog(*filename)
	if err != nil {
		panic(err)
	}
	mocker.Run()
}
