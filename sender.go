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
	Timestamp string `json: "Timestamp"`
	Log       string `json: "Log"`
}

const ESTimestampProp string = "\"Timestamp\": { \"type\": \"date\", \"format\": \"epoch_millis\" }"

func HttpCall(method, url string, bodyString string) {
	body := ioutil.NopCloser(bytes.NewBuffer([]byte(bodyString)))

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("[Sender] New request error, err[%s]", err)
	}
	client := &http.Client{Timeout: time.Duration(500) * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Sender] Put error, err[%s]", err)
	}
	resp_str, _ := httputil.DumpResponse(resp, true)
	log.Printf("[Sender] Response %s", resp_str)
}

func InitElasticSearchMapping(config *PeckTaskConfig) error {
	// Try init index mapping
	indexMapping := `{"mappings":{}}`
	uri := config.ESConfig.URL + "/" + config.ESConfig.Index
	log.Printf("[Sender] Init ElasticSearch mapping %s ", indexMapping)
	HttpCall(http.MethodPut, uri, indexMapping)

	// Try init Timestamp Field type
	propString := `{"properties":{"Timestamp":{"type":"date","format":"epoch_millis"}}}`
	uri = uri + "/_mappings/" + config.ESConfig.Type
	log.Printf("[Sender] Init ElasticSearch mapping %s ", propString)
	HttpCall(http.MethodPut, uri, propString)

	// Try init user fields type
	for _, v := range config.Fields {
		propS := `{"properties":{"` + v.Name + `":{"type":"` + v.Type + `"}}}`
		log.Printf("[Sender] Init ElasticSearch mapping %s ", propS)
		HttpCall(http.MethodPut, uri, propS)
	}
	return nil
}

func SendToElasticSearch(config *ElasticSearchConfig, fields map[string]string) {
	data := map[string]string{
		"Host":      GetHost(),
		"Timestamp": strconv.FormatInt(time.Now().UnixNano()/1000000, 10),
	}
	for k, v := range fields {
		data[k] = v
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
	} else {
		resp_str, _ := httputil.DumpResponse(resp, true)
		log.Printf("[Sender] Response %s", resp_str)
	}
}
