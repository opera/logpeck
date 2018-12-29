package logpeck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// InfluxDbConfig .
type InfluxDbConfig struct {
	Hosts    string `json:"Hosts"`
	Database string `json:"Database"`
}

// InfluxDbSender .
type InfluxDbSender struct {
	config        InfluxDbConfig
	mu            sync.Mutex
	lastIndexName string
	host          string
}

// NewInfluxDbSenderConfig .
func NewInfluxDbSenderConfig(jbyte []byte) (InfluxDbConfig, error) {
	influxDbConfig := InfluxDbConfig{}
	err := json.Unmarshal(jbyte, &influxDbConfig)
	if err != nil {
		return influxDbConfig, err
	}
	log.Infof("[NewInfluxDbSenderConfig]ElasticSearchConfig: %v", influxDbConfig)
	return influxDbConfig, nil
}

// NewInfluxDbSender .
func NewInfluxDbSender(senderConfig *SenderConfig) (*InfluxDbSender, error) {
	sender := InfluxDbSender{}
	config, ok := senderConfig.Config.(InfluxDbConfig)
	if !ok {
		return &sender, errors.New("New InfluxDbSender error ")
	}
	sender = InfluxDbSender{
		config: config,
	}

	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return &sender, errors.New("Get InfluxDbSender host error")
	}
	defer conn.Close()
	sender.host = strings.Split(conn.LocalAddr().String(), ":")[0]
	return &sender, nil
}

func (p *InfluxDbSender) toInfluxdbLine(fields map[string]interface{}) string {
	lines := ""
	timestamp := fields["timestamp"].(int64)

	for k, v := range fields {
		if k == "timestamp" {
			continue
		}
		aggregationResults := v.(map[string]float64)
		line := k + ",host=" + p.host + " "
		for aggregation, result := range aggregationResults {
			line += aggregation + "=" + strconv.FormatFloat(result, 'f', 3, 64) + ","
		}
		length := len(line)
		line = line[0:length-1] + " " + strconv.FormatInt(timestamp*1000000000, 10) + "\n"
		lines += line
		log.Infof("[toInfluxdbLine] line is %s", line)
	}
	return lines
}

// Start .
func (p *InfluxDbSender) Start() error {
	return nil
}

// Stop .
func (p *InfluxDbSender) Stop() error {
	return nil
}

// Send .
func (p *InfluxDbSender) Send(fields map[string]interface{}) {
	lines := p.toInfluxdbLine(fields)
	rawData := []byte(lines)
	body := ioutil.NopCloser(bytes.NewBuffer(rawData))
	uri := "http://" + p.config.Hosts + "/write?db=" + p.config.Database
	resp, err := http.Post(uri, "application/json", body)
	if err != nil {
		log.Infof("[InfluxDbSender.Sender] Post error, err[%s]", err)
	} else {
		respStr, _ := httputil.DumpResponse(resp, true)
		log.Infof("[InfluxDbSender.Sender] Response %s", respStr)
	}
	//p.measurments.MeasurmentRecall(fields)
}
