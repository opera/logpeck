package logpeck

import (
	"fmt"
	"testing"
)

func TestPeckTaskConfigUnmarshal(*testing.T) {
	var config PeckTaskConfig
	var configStr string

	configStr = `{
		"Name":"TestLog"
	}`
	if e := config.Unmarshal([]byte(configStr)); e != nil {
		panic(e)
	}

	configStr = `{
		"Nameeeee":"TestLog"
	}`
	if e := config.Unmarshal([]byte(configStr)); e == nil {
		panic("need field: Name")
	}

	configStr = `{
		"Name":"TestLog",
		"LogPath":"test.log",
		"ESConfig":{
			"Hosts":["127.0.0.1:9200","127.0.0.1:9201"],
			"Index":"TestLog",
			"Type":"Mock"
		}
	}`
	if e := config.Unmarshal([]byte(configStr)); e != nil {
		panic(e)
	}

	configStr = `{
		"Name":"TestLog",
		"LogPath":"test.log",
		"ESConfig":{
			"Hosts":["127.0.0.1:9200","127.0.0.1:9201"],
			"Index":"TestLog",
			"Type":"Mock",
			"Mapping":{
				"properties":"haha"
			}
		}
	}`
	if e := config.Unmarshal([]byte(configStr)); e != nil {
		panic(e)
	}
	fmt.Println(config)

	configStr = `{
		"Name":"TestLog",
		"LogPath":"test.log",
		"ESConfig":{
			"Hosts":["127.0.0.1:9200","127.0.0.1:9201"],
			"Index":"TestLog",
			"Type":"Mock",
			"Mapping":{
				"properties":"haha"
			}
		},
		"Delimiters": "",
		"FilterExpr":"mocklog hahaha|mocklog",
		"LogFormat": "json"
	}`
	if e := config.Unmarshal([]byte(configStr)); e != nil {
		panic(e)
	}
	fmt.Println(config)

	configStr = `{
		"Name":"TestLog",
		"LogPath":"test.log",
		"ESConfig":{
			"Hosts":["127.0.0.1:9200","127.0.0.1:9201"],
			"Index":"TestLog",
			"Type":"Mock",
			"Mapping":{
				"properties":"haha"
			}
		},
		"Fields":[
		{
			"Name": "DateText",
			"Value": "$1"
		},
		{
			"Name": "TS",
			"Value": "$6"
		}
		],
		"Delimiters": "",
		"FilterExpr":"mocklog hahaha|mocklog",
		"LogFormat": "json"
	}`
	if e := config.Unmarshal([]byte(configStr)); e != nil {
		//		panic(e)
	}
	fmt.Println(config)
}
