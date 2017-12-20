package logpeck

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"
)

/*
[{
"name":"module",
"tags":[
	"upstream",
	"downstream"
],
"aggr":["cnt","p99","avg"],
"target":"cost"
}]
*/

type InfluxDbConfig struct {
	Hosts       string                      `json:"hosts"`
	Interval    int64                       `json:"interval"`
	Name        string                      `json:"name"`
	DbName      string                      `json:"dbName"`
	Aggregators map[string]AggregatorConfig `json:"aggregators"`
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *InfluxDbSender {
	config := senderConfig.Config.(InfluxDbConfig)
	sender := InfluxDbSender{
		config: config,
		fields: fields,
	}
	return &sender
}

func toInfluxdbLine(fields map[string]interface{}) string {
	lines := ""
	timestamp := fields["timestamp"].(int64)
	host := ""
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		//return
	}
	defer conn.Close()
	host = strings.Split(conn.LocalAddr().String(), ":")[0]

	for k, v := range fields {
		if k == "timestamp" {
			continue
		}
		aggregationResults := v.(map[string]int64)
		lines = k + ",host=" + host + " "
		for aggregation, result := range aggregationResults {
			lines += aggregation + "=" + strconv.FormatInt(result, 10) + ","

		}
		length := len(lines)
		lines = lines[0:length-1] + " " + strconv.FormatInt(timestamp, 10) + "\n"
	}
	return lines
}

func (p *InfluxDbSender) Send(fields map[string]interface{}) {
	lines := toInfluxdbLine(fields)
	log.Infof("%s", lines)
	raw_data := []byte(lines)
	body := ioutil.NopCloser(bytes.NewBuffer(raw_data))
	uri := "http://" + p.config.Hosts + "/write?db=" + p.config.DbName
	resp, err := http.Post(uri, "application/json", body)
	if err != nil {
		log.Infof("[InfluxDbSender.Sender] Post error, err[%s]", err)
	} else {
		resp_str, _ := httputil.DumpResponse(resp, true)
		log.Debugf("[InfluxDbSender.Sender] Response %s", resp_str)
	}
	//p.measurments.MeasurmentRecall(fields)
}
