# Logpeck - An Interactive Log Collector

[![Documentation Status](https://img.shields.io/badge/中文文档-最新-brightgreen.svg)](README-cn.md)

## Features
 * Support plain text and json format log.
 * Full control with HTTP API.
 * Add/Remove/Start/Stop collection task freely.
 * Cooperate with [ElasticSearch](https://github.com/elastic/elasticsearch)/[Kibana](https://github.com/elastic/kibana) deeply.
 * Collection speed control.
 * Get collection status freely.
 * Web show/control conveniently([logpeck-web](https://github.com/opera/logpeck-web)).
 
## Build

`go build cmd/logpeckd/logpeckd.go`

## Getting Started

#### Requirements

Logpeck will post log data into an elasticsearch service, so there should be an elasticsearch service first. See [here](https://github.com/elastic/elasticsearch) for more information.

#### Launch logpeck service
 
`./logpeckd -config logpeckd.conf`

We can also use `supervisor` or other service management software to manage logpeck process.

#### Try collect a log

1. Add a new task first.

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

2. Start peck task.

```
curl -XPOST http://127.0.0.1:7117/peck_task/start -d {
  	"Name":"SystemLog"
}
```
```
Start Success
```

3. Stop peck task

```
curl -XPOST http://127.0.0.1:7117/peck_task/stop -d {
  	"Name":"SystemLog"
}
```

4. Remove peck task

```
curl -XPOST http://127.0.0.1:7117/peck_task/remove -d {
  	"Name":"SystemLog"
}
```

5. List peck tasks

```
curl -XPOST http://127.0.0.1:7117/peck_task/list
```

## Peck task configuration

## Http API
