# Logpeck - 交互式日志采集组件

[![Build Status](https://travis-ci.org/opera/logpeck.svg?branch=master)](https://travis-ci.org/opera/logpeck)
[![Documentation Status](https://img.shields.io/badge/English-Doc-brightgreen.svg)](README.md)

## 主要功能
 * 与[ElasticSearch](https://github.com/elastic/elasticsearch)/[Kibana](https://github.com/elastic/kibana)深度集成，方便使用其强大的存储、检索、可视化等功能
 * 支持普通文本及Json格式日志采集
 * 支持按表达式进行过滤、分段采集
 * 通过HTTP协议对日志收集进行控制和管理，方便集群进行集中管理
 * 可以随时对日志收集任务进行填加、暂停、更新、删除等操作
 * 容易支持Web端管理([logpeck-web](https://github.com/opera/logpeck-web))
 
## 构建

`go build cmd/logpeckd/logpeckd.go`

## 快速使用

以系统syslog为例进行日志采集。

#### 环境要求

Logpeck需要利用ElasticSearch进行数据存储和检索，使用Logpeck前，先保证有一个可用的ElasticSearch服务，详情[参见](https://github.com/elastic/elasticsearch)。

#### 启动
 
`./logpeckd -config logpeckd.conf`

也可以使用`supervisor`或其它管理软件对Logpeck进程进行管理。

#### 日志采集

1. 新增采集任务

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

2. 开始采集

```
curl -XPOST http://127.0.0.1:7117/peck_task/start -d {
  	"Name":"SystemLog"
}
```
```
Start Success
```

此时应该已经可以将`/var/log/syslog`中新增的日志写入配置好的ElasticSearch中。

3. 暂停采集

```
curl -XPOST http://127.0.0.1:7117/peck_task/stop -d {
  	"Name":"SystemLog"
}
```

4. 删除任务

```
curl -XPOST http://127.0.0.1:7117/peck_task/remove -d {
  	"Name":"SystemLog"
}
```

5. 列出所有采集任务

```
curl -XPOST http://127.0.0.1:7117/peck_task/list
```
