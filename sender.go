package logpeck

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type ElasticSearchSender struct {
	url string
}

func NewElasticSearchSender(url string) {

}

type ElasticSearchData struct {
	Host      string `json: "Host"`
	Log       string `json: "Log"`
	Timestamp string `json: "Timestamp"`
}

const ESTimestampProp string = "\"Timestamp\": { \"type\": \"date\", \"format\": \"epoch_millis\" }"

func HttpCall(method, url string, bodyRaw []byte) {
	body := ioutil.NopCloser(bytes.NewBuffer(bodyRaw))

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("[Sender] New request error, err[%s]", err)
	}
	//resp, err := http.Put(uri, "application/json", body)
	client := &http.Client{Timeout: time.Duration(500) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Sender] Put error, err[%s]", err)
	}
	resp_str, _ := httputil.DumpResponse(resp, true)
	log.Printf("[Sender] Response %s", resp_str)
}

func InitElasticSearchMapping(config *ElasticSearchConfig) {

	// Try init index mapping
	indexMapping := []byte(`{"mappings":{}}`)
	uri := config.URL + "/" + config.Index
	log.Printf("[Sender] Init ElasticSearch mapping %s ", indexMapping)
	HttpCall(http.MethodPut, uri, indexMapping)

	// Try init type mapping
	propString := []byte(`{"properties":{"Timestamp":{"type":"date","format":"epoch_millis"}}}`)
	uri = config.URL + "/" + config.Index + "/_mappings/" + config.Type
	log.Printf("[Sender] Init ElasticSearch mapping %s ", propString)
	HttpCall(http.MethodPut, uri, propString)
}

func SendToElasticSearch(config *ElasticSearchConfig, content string) {
	data := ElasticSearchData{
		Host:      GetHost(),
		Log:       content,
		Timestamp: strconv.FormatInt(time.Now().UnixNano()/1000000, 10),
	}
	log.Println(data)
	raw_data, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	uri := config.URL + "/" + config.Index + "/" + config.Type
	log.Printf("[Sender] Post ElasticSearch %s content [%s] ", uri, raw_data)
	body := ioutil.NopCloser(bytes.NewBuffer(raw_data))
	resp, err := http.Post(uri, "application/json", body)
	if err != nil {
		log.Printf("[Sender] Post error, err[%s]", err)
	}
	resp_str, _ := httputil.DumpResponse(resp, true)
	log.Printf("[Sender] Response %s", resp_str)
}
