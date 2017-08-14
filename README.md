# Logpeck - A Simple, RESTful Log Collector

[![Build Status](https://travis-ci.org/opera/logpeck.svg?branch=master)](https://travis-ci.org/opera/logpeck)
[![Documentation Status](https://img.shields.io/badge/中文文档-最新-brightgreen.svg)](README-cn.md)

## Objectives
Logpeck aims to be an easy-to-use module that collecting contents from log file and posting into [ElasticSearch](https://github.com/elastic/elasticsearch). We want to control collection tasks remotely with HTTP API (**NONE configuration file**).

## Build

`go build cmd/logpeckd/logpeckd.go`

## Getting Started

#### Launch logpeck service
 
`./logpeckd -config logpeckd.conf`

We can also use `supervisor` or other service management software to manage logpeck process.

#### Try collect a log

1. Add a new task first. (Want more task config, filter, json, long, etc.? see [here](doc/task_config.md).)

```
curl -XPOST http://127.0.0.1:7117/peck_task/add -d {
  	"Name":"SystemLog",
	"LogPath":"/var/log/syslog",
	"ESConfig":{
	  	"Hosts":["127.0.0.1:9200"],
		"Index":"syslog",
		"Type":"raw"
	}
}
```
```
Add Success
```

2. Start task.

```
curl -XPOST http://127.0.0.1:7117/peck_task/start -d {
  	"Name":"SystemLog"
}
```
```
Start Success
```

3. Stop task

```
curl -XPOST http://127.0.0.1:7117/peck_task/stop -d {
  	"Name":"SystemLog"
}
```

4. Remove task

```
curl -XPOST http://127.0.0.1:7117/peck_task/remove -d {
  	"Name":"SystemLog"
}
```

5. List tasks

```
curl -XPOST http://127.0.0.1:7117/peck_task/list
```

## Documentation

 * [Peck Task Configuration](doc/task_config.md)
 * [Frequently Asked Questions](doc/FAQ.md)
