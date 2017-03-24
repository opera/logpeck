package logpeck

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func NewAddTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddTaskHandler")
		r_str, _ := httputil.DumpRequest(r, false)
		log.Printf("Request len[%d], body[%s]", len(r_str), r_str)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("AddTaskHandler Success\n"))
	}
}

func NewUpdateTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("UpdateTaskHandler")
		r_str, _ := httputil.DumpRequest(r, false)
		log.Printf("Request len[%d], body[%s]", len(r_str), r_str)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("UpdateTaskHandler Success\n"))
	}
}

func NewPauseTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("PauseTaskHandler")
		r_str, _ := httputil.DumpRequest(r, false)
		log.Printf("Request len[%d], body[%s]", len(r_str), r_str)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PauseTaskHandler Success\n"))
	}
}

func NewRemoveTaskHandler(pecker *Pecker, db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("RemoveTaskHandler")
		r_str, _ := httputil.DumpRequest(r, false)
		log.Printf("Request len[%d], body[%s]", len(r_str), r_str)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("RemoveTaskHandler Success\n"))
	}
}
