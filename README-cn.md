# Logpeck - 简洁灵活的日志收集组件

[![Build Status](https://travis-ci.org/opera/logpeck.svg?branch=master)](https://travis-ci.org/opera/logpeck)
[![Documentation Status](https://img.shields.io/badge/English-Doc-brightgreen.svg)](README.md)

## 功能
Logpeck尝试用最简洁灵活的方式收集并解析日志文件，将数据推送至不同的存储系统中去（比如[ElasticSearch](https://github.com/elastic/elasticsearch), [Influxdb](https://github.com/influxdata/influxdb), [Kafka](https://github.com/apache/kafka)）。

 * Logpeck通过HTTP API的方式进行任务管理、更新，脱离配置文件。
 * 支持多种方法解析日志
  * 按列分隔－－简单有效，节省资源
  * 内嵌Lua－－功能强大，可处理任意格式日志
  * Json－－Json格式数据最有效
 * 方便支持Web UI, 默认提供[logpeck-kibana-plugin](https://github.com/opera/logpeck-kibana-plugin)

## 使用

### 安装与启动
#### 安装包

 * 仅提供linux，其它系统需自己编译
 * 下载安装包 [logpeck_0.5.0.deb](https://github.com/opera/resources/blob/master/logpeck/releases/logpeck_0.5.0.deb)
 * 安装： `sudo dpkg -i logpeck_0.5.0.deb`
 * 启动： `sudo service logpeck start` (如果支持`supervisor`，可使用 `sudo supervisorctl update`) 

#### 源码编译

 * 下载源代码: [Release page v0.5.0](https://github.com/opera/logpeck/releases/tag/0.5.0)
 * 编译： `go build cmd/logpeckd/logpeckd.go`
 * 启动： `./logpeckd -config logpeckd.conf`

### 可视化界面

[**logpeck-kibana-plugin**](https://github.com/opera/logpeck-kibana-plugin) 是默认提供的Logpeck可视化界面，通过此插件可以使用Logpeck的全部功能，并提供任务及集群管理功能，推荐使用。

## 文档

 * [HTTP API协议](doc/task_config.md)
 * [FAQ](doc/FAQ.md)
 
## 项目依赖

 * [BurntSushi/toml](https://github.com/BurntSushi/toml)
 * [Sirupsen/logrus](https://github.com/Sirupsen/logrus)
 * [bitly/go-simplejson](https://github.com/bitly/go-simplejson)
 * [yuin/gopher-lua](https://github.com/yuin/gopher-lua)
 * [boltdb/bolt](https://github.com/boltdb/bolt)
 * [go-zoo/bone](https://github.com/go-zoo/bone)
 * [hpcloud/tail](https://github.com/hpcloud/tail)
 * [Shopify/sarama](https://github.com/Shopify/sarama)
 
由衷感谢
 
## 讨论

有任何建议或问题, 通过[issue](https://github.com/opera/logpeck/issues/new)反馈。
