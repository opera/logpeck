package logpeck

import (
	"strconv"
	"testing"

	log "github.com/Sirupsen/logrus"
)

/*
func TestStartSend(*testing.T) {
	log.Infof("[aggregator_test] TestStartSend")

	test := AggregatorOption{
		PreMeasurment: "Test",
		Measurment:    "aaa",
		Tags:          []string{"cost"},
		Aggregations:  []string{"cnt"},
		Target:        "cost",
		Timestamp:     "time",
	}
	var options []AggregatorOption
	options = append(options, test)
	aggregatorConfig := AggregatorConfig{
		Enable:   true,
		Interval: int64(30),
		Options:  options,
	}
	aggregator := NewAggregator(&aggregatorConfig)
	aggregator.recordTime = 29
}
*/

func TestRecord(*testing.T) {
	test := AggregatorOption{
		PreMeasurment: "Test",
		Measurment:    "aaa",
		Tags:          []string{"upstream"},
		Aggregations:  []string{"cnt"},
		Target:        "cost",
		Timestamp:     "time",
	}
	var options []AggregatorOption
	options = append(options, test)
	aggregatorConfig := AggregatorConfig{
		Enable:   true,
		Interval: int64(30),
		Options:  options,
	}
	aggregator := NewAggregator(&aggregatorConfig)

	fields := make(map[string]interface{})
	fields["aaa"] = "getTest"
	fields["upstream"] = "127.0.0.1"
	fields["cost"] = "2"
	fields["time"] = "15"
	aggregator.Record(fields)
	if aggregator.buckets["Test_aaa_cost"]["Test_getTest_cost,upstream=127.0.0.1"][0] != 2 {
		panic(aggregator)
	}
	aggregator.Record(fields)
	if aggregator.buckets["Test_aaa_cost"]["Test_getTest_cost,upstream=127.0.0.1"][0]+aggregator.buckets["Test_aaa_cost"]["Test_getTest_cost,upstream=127.0.0.1"][1] != 4 {
		panic(aggregator)
	}
}

func TestDump(*testing.T) {
	test := AggregatorOption{
		PreMeasurment: "Test",
		Measurment:    "aaa",
		Tags:          []string{"upstream"},
		Aggregations:  []string{"cnt", "avg", "p99", "p50"},
		Target:        "cost",
		Timestamp:     "time",
	}
	var options []AggregatorOption
	options = append(options, test)
	aggregatorConfig := AggregatorConfig{
		Enable:   true,
		Interval: int64(30),
		Options:  options,
	}
	aggregator := NewAggregator(&aggregatorConfig)

	fields := make(map[string]interface{})
	fields["aaa"] = "getTest"
	fields["upstream"] = "127.0.0.1"
	fields["cost"] = "0"
	fields["time"] = "15"
	for i := 9; i >= 0; i-- {
		fields["cost"] = strconv.Itoa(i)
		aggregator.Record(fields)
	}
	dump := aggregator.Dump()
	log.Infof("%v", dump)
	a := dump["Test_getTest_cost,upstream=127.0.0.1"].(map[string]float64)
	if a["cnt"] != 10 {
		panic(a)
	}
	if a["avg"] != 4.5 {
		log.Panicf("%#v", a)
	}
	if a["p99"] != 8 {
		panic(a)
	}
	if a["p50"] != 4 {
		panic(a)
	}
}
