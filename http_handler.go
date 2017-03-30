package logpeck

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

func logRequest(r *http.Request, prefix string) {
	r_str, _ := httputil.DumpRequest(r, true)
	log.Printf("[Handler] [%s] req_len[%d] req[%s]", prefix, len(r_str), r_str)
}

func NewAddTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "AddTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(raw, &config)
		if err != nil {
			log.Printf("[Handler] Parse PeckTaskConfig error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request\n"))
			return
		}

		err = pecker.AddPeckTask(&config)
		if err != nil {
			log.Printf("[Handler] AddTaskConfig error, %s", err)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Add failed, " + err.Error() + "\n"))
			return
		}
		log.Printf("[Handler] Add PeckTask Success: %s", raw)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
		return
	}
}

func NewUpdateTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "UpdateTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(raw, &config)
		if err != nil {
			log.Printf("[Handler] Parse PeckTaskConfig error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request\n"))
			return
		}

		err = pecker.UpdatePeckTask(&config)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Update failed, " + err.Error() + "\n"))
			return
		}

		err = db.SaveConfig(&config)
		if err != nil {
			log.Printf("[Handler] UpdateTaskConfig error, save config error, %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("UpdateTaskHandler Success\n"))
	}
}

func NewStartTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "StartTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(raw, &config)
		if err != nil {
			log.Printf("[Handler] Start PeckTask error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request, " + err.Error() + "\n"))
			return
		}

		err = pecker.StartPeckTask(&config)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Start failed, " + err.Error() + "\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("StartTask Success\n"))
	}
}

func NewPauseTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "PauseTaskHandler")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PauseTaskHandler Success\n"))
	}
}

func NewRemoveTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "RemoveTaskHandler")
		defer r.Body.Close()

		var config PeckTaskConfig
		raw, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(raw, &config)
		if err != nil {
			log.Printf("[Handler] Remove PeckTask error, %s", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request, " + err.Error() + "\n"))
			return
		}

		err = pecker.RemovePeckTask(&config)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Remove PeckTask failed, " + err.Error() + "\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\n"))
	}
}
