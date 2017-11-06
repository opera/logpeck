package logpeck

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func logRequest(r *http.Request, prefix string) {
	r_str, _ := httputil.DumpRequest(r, true)
	log.Infof("[Handler] [%s] req_len[%d] req[%s]", prefix, len(r_str), r_str)
}

func NewAddTaskHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "AddTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := config.Unmarshal(raw)
		if err != nil {
			log.Infof("[Handler] Parse PeckTaskConfig error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, %s in %v", err, string(raw[:]))))
			return
		}

		err = pecker.AddPeckTask(&config, nil)
		if err != nil {
			log.Infof("[Handler] AddTaskConfig error, %s", err)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Add failed, " + err.Error()))
			return
		}
		log.Infof("[Handler] Add Success: %s", raw)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Add Success"))
		return
	}
}

func NewUpdateTaskHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "UpdateTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := config.Unmarshal(raw)
		if err != nil {
			log.Infof("[Handler] Parse PeckTaskConfig error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Bad Request, %s in %v", err, string(raw[:]))))
			return
		}

		err = pecker.UpdatePeckTask(&config)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Update failed, " + err.Error()))
			return
		}

		if err != nil {
			log.Infof("[Handler] UpdateTaskConfig error, save config error, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		log.Infof("[Handler] Update Success: %s", raw)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update Success"))
	}
}

func NewStartTaskHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "StartTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := config.Unmarshal(raw)
		if err != nil {
			log.Infof("[Handler] Start PeckTask error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request, " + err.Error()))
			return
		}

		err = pecker.StartPeckTask(&config)
		if err != nil {
			log.Infof("[Handler] Start PeckTask error, %s", err.Error())
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Start failed, " + err.Error()))
			return
		}
		log.Infof("[Handler] Start Success: %s", raw)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Start Success"))
	}
}

func NewStopTaskHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "StopTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := config.Unmarshal(raw)
		if err != nil {
			log.Infof("[Handler] Stop PeckTask error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request, " + err.Error()))
			return
		}

		err = pecker.StopPeckTask(&config)
		if err != nil {
			log.Infof("[Handler] Stop PeckTask error, %s", err.Error())
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Stop failed, " + err.Error()))
			return
		}
		log.Infof("[Handler] Stop Success: %s", raw)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Stop Success"))
	}
}

func NewRemoveTaskHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "RemoveTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := config.Unmarshal(raw)
		if err != nil {
			log.Infof("[Handler] Remove PeckTask error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request, " + err.Error()))
			return
		}

		err = pecker.RemovePeckTask(&config)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Remove PeckTask failed, " + err.Error()))
			return
		}
		log.Infof("[Handler] Remove Success: %s", raw)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Remove Success"))
	}
}

func NewListTaskHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "ListTaskHandler")
		defer r.Body.Close()

		configs, err := pecker.ListPeckTask()
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("List PeckTask failed, " + err.Error()))
			return
		}
		jsonStr, jErr := json.Marshal(configs)
		if jErr != nil {
			panic(jErr)
		}
		log.Infof("[Handler] List Success: %s", jsonStr)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonStr))
	}
}

func NewListStatsHandler(pecker *Pecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "ListStatusHandler")
		defer r.Body.Close()

		stats, err := pecker.ListTaskStats()
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("List TaskStatus failed, " + err.Error()))
			return
		}
		jsonStr, jErr := json.Marshal(stats)
		if jErr != nil {
			panic(jErr)
		}
		log.Infof("[Handler] List Success: %s", jsonStr)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonStr))
	}
}
