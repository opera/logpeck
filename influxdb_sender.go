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
	"time"
)

type InfluxDbConfig struct {
	Hosts       string                      `json:"Hosts"`
	DBName      string                      `json:"DBName"`
	Interval    int64                       `json:"Interval"`
	FieldsKey   string                      `json:"FieldsKey"`
	Aggregators map[string]AggregatorConfig `json:"Aggregators"`
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	taskName      string
	mu            sync.Mutex
	lastIndexName string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField, taskName string) *InfluxDbSender {
	config := senderConfig.Config.(InfluxDbConfig)
	sender := InfluxDbSender{
		config:   config,
		fields:   fields,
		taskName: taskName,
	}
	return &sender
}

func toInfluxdbLine(fields map[string]interface{}, taskName string) string {
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
		key := strings.Split(k, " ")
		lines += taskName + "_" + key[0] + ",host=" + host + " "
		for aggregation, result := range aggregationResults {
			lines += key[1] + aggregation + "=" + strconv.FormatInt(result, 10) + ","
		}
		length := len(lines)
		lines = lines[0:length-1] + " " + strconv.FormatInt(timestamp*1000000000, 10) + "\n"
	}
	return lines
}

func (p *InfluxDbSender) Send(fields map[string]interface{}) {
	lines := toInfluxdbLine(fields, p.taskName)
	log.Infof("[InfluxDbSender.Sender] timestamp is %v", time.Now())
	log.Infof("[InfluxDbSender.Sender] lines is %s", lines)
	raw_data := []byte(lines)
	body := ioutil.NopCloser(bytes.NewBuffer(raw_data))
	uri := "http://" + p.config.Hosts + "/write?db=" + p.config.DBName
	resp, err := http.Post(uri, "application/json", body)
	if err != nil {
		log.Infof("[InfluxDbSender.Sender] Post error, err[%s]", err)
	} else {
		resp_str, _ := httputil.DumpResponse(resp, true)
		log.Infof("[InfluxDbSender.Sender] Response %s", resp_str)
	}
	//p.measurments.MeasurmentRecall(fields)
}
