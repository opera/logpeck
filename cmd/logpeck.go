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

func restorePeckTasks(pecker *logpeck.Pecker, db *logpeck.DB) {
	defer logpeck.LogExecTime(time.Now(), "Restore PeckTaskConfig")
	configs, err := db.GetAllConfigs()
	if err != nil {
		panic(err)
	}
	for i, config := range configs {
		pecker.AddPeckTask(&config)
		log.Printf("Restore PeckTask[%d] : %s", i, config)
	}
}

func main() {
	configFile := flag.String("config", "./logpeck.conf", "Config file path")
	flag.Parse()

	logpeck.InitConfig(configFile)

	pecker, err := logpeck.NewPecker()
	if err != nil {
		panic(err)
	}

	err = logpeck.OpenDB(logpeck.Config.DatabaseFile)
	if err != nil {
		panic(err)
	}
	db := logpeck.GetDBHandler()
	defer db.Close()

	restorePeckTasks(pecker, db)

	mux := bone.New()
	mux.Post("/peck_task/add", logpeck.NewAddTaskHandler(pecker, db))
	mux.Post("/peck_task/update", logpeck.NewUpdateTaskHandler(pecker, db))
	mux.Post("/peck_task/pause", logpeck.NewPauseTaskHandler(pecker, db))
	mux.Post("/peck_task/remove", logpeck.NewRemoveTaskHandler(pecker, db))
	//	mux.Get("/pecker_stat", http.HandlerFunc(handler.Get))
	//	mux.Post("/peck_task/remove", http.HandlerFunc(handler.NewsRedirectHandler))

	address := fmt.Sprintf(":%d", logpeck.Config.Port)

	log.Printf("Logpeck start serving on port %d ...\n", logpeck.Config.Port)
	s := &http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}
	s.ListenAndServe()
}
