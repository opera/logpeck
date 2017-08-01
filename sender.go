package logpeck

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

type ElasticSearchSender struct {
	config ElasticSearchConfig
}

func NewElasticSearchSender(config *ElasticSearchConfig) *ElasticSearchSender {
	return &ElasticSearchSender{
		config: *config,
	}
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
	} else {
		resp_str, _ := httputil.DumpResponse(resp, true)
		log.Printf("[Sender] Response %s", resp_str)
	}
}

func GetIndexName(prototype string) string {
	l, r := "%{+", "}"
	if !strings.Contains(prototype, l) || !strings.Contains(prototype, r) {
		return prototype
	}
	indexName := ""
	lIndex := strings.Index(prototype, l)
	rIndex := strings.Index(prototype, r)
	format := prototype[lIndex+len(l) : rIndex]
	timeStr := time.Now().Format(format)

	return indexName + prototype[:lIndex] + timeStr + prototype[rIndex+1:]
}

func (p *ElasticSearchSender) Init(taskConfig *PeckTaskConfig) error {
	// Try init index mapping
	indexMapping := `{"mappings":{}}`
	host, err := SelectRandom(p.config.Hosts)
	if err != nil {
		return err
	}
	uri := "http://" + host + "/" + GetIndexName(p.config.Index)
	log.Printf("[Sender] Init ElasticSearch mapping %s ", indexMapping)
	HttpCall(http.MethodPut, uri, indexMapping)

	// Try init Timestamp Field type
	propString := `{"properties":{"Timestamp":{"type":"date","format":"epoch_millis"}}}`
	uri = uri + "/_mappings/" + p.config.Type
	log.Printf("[Sender] Init ElasticSearch mapping %s ", propString)
	HttpCall(http.MethodPut, uri, propString)

	// Try init user fields type
	for _, v := range taskConfig.Fields {
		propS := `{"properties":{"` + v.Name + `":{"type":"` + v.Type + `"}}}`
		log.Printf("[Sender] Init ElasticSearch mapping %s ", propS)
		HttpCall(http.MethodPut, uri, propS)
	}
	return nil
}

func (p *ElasticSearchSender) Send(fields map[string]interface{}) {
	data := map[string]interface{}{
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
	host, err := SelectRandom(p.config.Hosts)
	if err != nil {
		log.Printf("[Sender] ElasticSearch Host error [%v] ", err)
		return
	}
	uri := "http://" + host + "/" + GetIndexName(p.config.Index) + "/" + p.config.Type
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
