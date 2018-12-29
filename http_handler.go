package logpeck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-zoo/bone"
)

func logRequest(r *http.Request, prefix string) {
	str, _ := httputil.DumpRequest(r, true)
	log.Infof("[Handler] [%s] req_len[%d] req[%s]", prefix, len(str), str)
}

// NewAddTaskHandler .
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

// NewUpdateTaskHandler .
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

// NewStartTaskHandler .
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

// NewStopTaskHandler .
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

// NewRemoveTaskHandler .
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

// NewListTaskHandler .
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
		stats, err := pecker.ListTaskStats()
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("List PeckTask failed, " + err.Error()))
			return
		}
		res := make(map[string]interface{})
		res["configs"] = configs
		res["stats"] = stats
		jsonStr, jErr := json.Marshal(res)
		if jErr != nil {
			panic(jErr)
		}
		log.Infof("[Handler] List Success: %s", jsonStr)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonStr))
	}
}

// NewTestTaskHandler .
func NewTestTaskHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "TestTaskHandler")
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

		results, err := TestPeckTask(&config)
		if err != nil {
			log.Infof("[Handler] TestTaskConfig error, %s", err)
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("test failed, " + err.Error()))
			return
		}
		jsonStr, jErr := json.Marshal(results)
		if jErr != nil {
			panic(jErr)
		}
		log.Infof("[Handler] Test Success: %s", jsonStr)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
		return
	}
}

// NewListPathHandler .
func NewListPathHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "ListPathHandler")

		path := ""
		if pathArr := bone.GetQuery(r, "path"); len(pathArr) > 0 {
			path = pathArr[0]
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Need path"))
			return
		}

		index := strings.LastIndex(path, "/")
		if index == -1 {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Path should be absolute"))
		}
		dir := path[0 : index+1]
		filePrefix := path[index+1:]

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}

		var names []string
		for _, file := range files {
			if name := file.Name(); strings.HasPrefix(name, filePrefix) {
				names = append(names, strings.TrimPrefix(name, filePrefix))
				if len(names) >= 20 {
					break
				}
			}
		}

		jsonStr, err := json.Marshal(names)
		if err != nil {
			panic(err)
		}
		log.Infof("[Handler] List path Success: %s", jsonStr)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonStr))
	}
}

// NewVersionHandler .
func NewVersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, "VersionHandler")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(VersionString))
	}
}
