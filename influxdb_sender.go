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

type InfluxDbConfig struct {
	Hosts             string             `json:"Hosts"`
	DBName            string             `json:"DBName"`
	Interval          int64              `json:"Interval"`
	AggregatorConfigs []AggregatorConfig `json:"AggregatorConfigs"`
}

type InfluxDbSender struct {
	config        InfluxDbConfig
	fields        []PeckField
	mu            sync.Mutex
	lastIndexName string
	host          string
}

func NewInfluxDbSender(senderConfig *SenderConfig, fields []PeckField) *InfluxDbSender {
	config := senderConfig.Config.(InfluxDbConfig)

	sender := InfluxDbSender{
		config: config,
		fields: fields,
	}

	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer conn.Close()
	sender.host = strings.Split(conn.LocalAddr().String(), ":")[0]
	return &sender
}

func (p *InfluxDbSender) toInfluxdbLine(fields map[string]interface{}) string {
	lines := ""
	timestamp := fields["timestamp"].(int64)

	for k, v := range fields {
		if k == "timestamp" {
			continue
		}
		aggregationResults := v.(map[string]int64)
		line := k + ",host=" + p.host + " "
		for aggregation, result := range aggregationResults {
			line += aggregation + "=" + strconv.FormatInt(result, 10) + ","
		}
		length := len(line)
		line = line[0:length-1] + " " + strconv.FormatInt(timestamp*1000000000, 10) + "\n"
		lines += line
		log.Infof("[toInfluxdbLine] line is %s", line)
	}
	return lines
}

func (p *InfluxDbSender) Start() error {
	return nil
}

func (p *InfluxDbSender) Stop() error {
	return nil
}

func (p *InfluxDbSender) Send(fields map[string]interface{}) {
	lines := p.toInfluxdbLine(fields)
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
