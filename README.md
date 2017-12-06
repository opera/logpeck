# Logpeck - A Simple, RESTful Log Collector

[![Build Status](https://travis-ci.org/opera/logpeck.svg?branch=master)](https://travis-ci.org/opera/logpeck)
[![Documentation Status](https://img.shields.io/badge/中文文档-最新-brightgreen.svg)](README-cn.md)

## Objectives
Logpeck aims to be an easy-to-use module that parsing and collecting contents from log file and posting into [ElasticSearch](https://github.com/elastic/elasticsearch). We want to control collection tasks remotely with HTTP API (**NONE configuration file**).

## Getting Started

### Build & Launch

`go build cmd/logpeckd/logpeckd.go`

`./logpeckd -config logpeckd.conf`

We can also use `supervisor` or other service management software to manage logpeck process.

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

 * [BurntSushi/toml](github.com/BurntSushi/toml): configuration management
 * [Sirupsen/logrus](github.com/Sirupsen/logrus): logging
 * [bitly/go-simplejson](github.com/bitly/go-simplejson): json parser
 * [boltdb/bolt](github.com/boltdb/bolt): local storage
 * [go-zoo/bone](github.com/go-zoo/bone): http multiplexer
 * [hpcloud/tail](github.com/hpcloud/tail): watching log file
 
 Saulte to all these excellent projects.
 
## Discussion

Any suggestions or questions, please [create an issue](https://github.com/opera/logpeck/issues/new) to feedback.
