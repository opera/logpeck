package logpeck

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

type ElasticSearchSender struct {
	url string
}

func NewElasticSearchSender(url string) {

}

type ElasticSearchData struct {
	Host string
	Log  string
}

func SendToElasticSearch(url, index, ty, content string) {
	data := ElasticSearchData{
		Host: GetHost(),
		Log:  content,
	}
	log.Println(data)
	raw_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	uri := url + "/" + index + "/" + ty
	log.Printf("[Sender] Post ElasticSearch %s content [%s] ", url, raw_data)
	body := ioutil.NopCloser(bytes.NewBuffer(raw_data))
	resp, err := http.Post(uri, "application/json", body)
	if err != nil {
		log.Printf("[Sender] Post error, err[%s]", err)
	}
	resp_str, _ := httputil.DumpResponse(resp, true)
	log.Printf("[Sender] Response %s", resp_str)
}
