package main

import (
	"flag"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/opera/logpeck"
	"net/http"
	"time"
)

func main() {
	configFile := flag.String("config", "./logpeck.conf", "Config file path")
	flag.Parse()

	logpeck.InitConfig(configFile)

	pecker, err := logpeck.NewPecker()
	if err != nil {
		panic(err)
	}
	fmt.Println(pecker)

	err = logpeck.OpenDB(logpeck.Config.DatabaseFile)
	if err != nil {
		panic(err)
	}
	db := logpeck.GetDBHandler()
	defer db.Close()

	mux := bone.New()
	//	mux.Get("/pecker_stat", http.HandlerFunc(handler.Get))
	//	mux.Post("/peck_task/add", http.HandlerFunc(handler.NewsRedirectHandler))
	//	mux.Post("/peck_task/remove", http.HandlerFunc(handler.NewsRedirectHandler))

	address := fmt.Sprintf(":%d", logpeck.Config.Port)

	fmt.Printf("Logpeck start serving on %s .\n", address)
	s := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	s.ListenAndServe()
}
