package main

import (
	"fmt"
	"github.com/opera/logpeck"
)

func main() {
	fmt.Println("hello logpeck")
	pecker, err := logpeck.NewPecker()
	if err != nil {
		panic(err)
	}
	fmt.Println(pecker)
}
