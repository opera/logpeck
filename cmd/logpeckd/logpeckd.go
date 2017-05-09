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
	configFile := flag.String("config", "./logpeckd.conf", "Config file path")
	flag.Parse()

	logpeck.InitConfig(configFile)
	log.Printf("[LogPeckD] LogPeck(%s) Start %+v", logpeck.VersionString, logpeck.Config)

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
	mux.Post("/peck_task/add", logpeck.NewAddTaskHandler(pecker))
	mux.Post("/peck_task/update", logpeck.NewUpdateTaskHandler(pecker))
	mux.Post("/peck_task/start", logpeck.NewStartTaskHandler(pecker))
	mux.Post("/peck_task/stop", logpeck.NewStopTaskHandler(pecker))
	mux.Post("/peck_task/remove", logpeck.NewRemoveTaskHandler(pecker))
	mux.Post("/peck_task/list", logpeck.NewListTaskHandler(pecker))

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
