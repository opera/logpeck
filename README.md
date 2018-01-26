# Logpeck - A Simple, RESTful Log Collector

[![Build Status](https://travis-ci.org/opera/logpeck.svg?branch=master)](https://travis-ci.org/opera/logpeck)
[![Documentation Status](https://img.shields.io/badge/中文文档-最新-brightgreen.svg)](README-cn.md)

## Objectives
Logpeck aims to be an easy-to-use module that parsing and collecting contents from log file and posting into specific storage system, such as [ElasticSearch](https://github.com/elastic/elasticsearch), [Influxdb](https://github.com/influxdata/influxdb), [Kafka](https://github.com/apache/kafka). We want to control collection tasks remotely with HTTP API (**NONE configuration file**).

## Getting Started

### Installation
#### From Binary (linux only)

 * Download installation package [logpeck_0.3.0.deb](https://github.com/opera/resources/blob/master/logpeck/releases/logpeck_0.3.0.deb)
 * Run `sudo dpkg -i logpeck_0.3.0.deb`
 * Run `sudo service logpeck start` (or `sudo supervisorctl update` if `supervisor` is avalible) 

#### From Source Code

 * Download source code: [Release page v0.3.0](https://github.com/opera/logpeck/releases/tag/0.3.0)
 * Build: `go build cmd/logpeckd/logpeckd.go`
 * Launch: `./logpeckd -config logpeckd.conf`
 * We can also use `supervisor` or other service management software to manage logpeck process.

### Web UI

We highly recommend to install [**logpeck-kibana-plugin**](https://github.com/opera/logpeck-kibana-plugin) into [Kibana](https://github.com/elastic/kibana). With this plugin, we can control all machines and collection tasks conveniently. At the same time, we can take advantage of powerful searching and visualization features of Kibana.

<p float="left">
  <img src="https://github.com/opera/resources/blob/master/logpeck/1.png" width="400" />
  <img src="https://github.com/opera/resources/blob/master/logpeck/2.png" width="400" /> 
</p>

### RESTful API

We can also control collection tasks with RESTful API. [See more](doc/restful.md)

## Documentation

 * [Peck Task Configuration](doc/task_config.md)
 * [Frequently Asked Questions](doc/FAQ.md)
 
## Dependencies

 * [BurntSushi/toml](https://github.com/BurntSushi/toml): configuration management
 * [Sirupsen/logrus](https://github.com/Sirupsen/logrus): logging
 * [bitly/go-simplejson](https://github.com/bitly/go-simplejson): json parser
 * [yuin/gopher-lua](https://github.com/yuin/gopher-lua): lua virtual machine
 * [boltdb/bolt](https://github.com/boltdb/bolt): local storage
 * [go-zoo/bone](https://github.com/go-zoo/bone): http multiplexer
 * [hpcloud/tail](https://github.com/hpcloud/tail): watching log file
 * [Shopify/sarama](https://github.com/Shopify/sarama): kafka client
 
 Saulte to all these excellent projects.
 
## Discussion

Any suggestions or questions, please [create an issue](https://github.com/opera/logpeck/issues/new) to feedback.
