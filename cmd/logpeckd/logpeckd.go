package main

import (
	"flag"
	"fmt"
	"github.com/go-zoo/bone"
	"github.com/opera/logpeck"
	"log"
	"net/http"
	"time"
)

func main() {
	configFile := flag.String("config", "./logpeck.conf", "Config file path")
	flag.Parse()

	logpeck.InitConfig(configFile)
	log.Printf("[LogPeckD] Try create a new logpeck: %s", logpeck.Config)

	err := logpeck.OpenDB(logpeck.Config.DatabaseFile)
	if err != nil {
		panic(err)
	}
	db := logpeck.GetDBHandler()
	defer db.Close()

	pecker, p_err := logpeck.NewPecker(db)
	if p_err != nil {
		panic(p_err)
	}
	pecker.Start()

	mux := bone.New()
	mux.Post("/peck_task/add", logpeck.NewAddTaskHandler(pecker, db))
	mux.Post("/peck_task/update", logpeck.NewUpdateTaskHandler(pecker, db))
	mux.Post("/peck_task/start", logpeck.NewStartTaskHandler(pecker, db))
	mux.Post("/peck_task/pause", logpeck.NewPauseTaskHandler(pecker, db))
	mux.Post("/peck_task/remove", logpeck.NewRemoveTaskHandler(pecker, db))

	//	mux.Get("/pecker_stat", http.HandlerFunc(handler.Get))

	log.Printf("[LogPeckD] Logpeck start serving on port %d ...\n", logpeck.Config.Port)
	address := fmt.Sprintf(":%d", logpeck.Config.Port)
	s := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	s.ListenAndServe()
}
