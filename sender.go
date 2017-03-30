package logpeck

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	log.Printf("Post ElasticSearch %s content [%s] ", url, raw_data)
	body := ioutil.NopCloser(bytes.NewBuffer(raw_data))
	resp, err := http.Post(uri, "application/json", body)
	if err != nil {
		log.Printf("Post error, err[%s], resp[%s]", err, resp)
	}
	log.Printf("Response %s err %s", resp, err)
}
