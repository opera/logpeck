# Peck Task Configuration

Logpeck use a json string to define a task. A simplest config is as follows.

```
{
  "Name": "HttpServer",
  "LogPath": "/data/log/http_server.log",
  "Keywords": "Performace",
  "Sender": {
    "Name": "Elasticsearch",
    "Config": {
      "Index": "http_server",
      "Type": "perf",
      "Hosts": [
        "10.0.0.11:9200",
        "10.0.0.12:9200",
        "10.0.0.13:9200"
      ]
    }
  },
  "Extractor": {
    "Name": "text",
    "Config": {
      "Fields": [
        {
          "Name": "module",
          "Value": "$6"
        },
        {
          "Name": "server",
          "Value": "$7"
        }
      ]
    }
  }
}
```

#### Name

A unique identification of a peck task. Logpeck use Name to control the specific task such as start/stop/update/remove/etc.

#### LogPath

The log file path to be pecked.

File should be accessible. 

If the file is not exist, task will still keep pecking. If the file is rotated, task will peck the new file named "LogPath".

#### ESConfig

 1. Hosts: ElasticSearch service hosts. logpeck will select randomly from this host list.
 2. Index: ElasticSearch index name.
 3. Index: ElasticSearch type name.

## Optional Configuration

#### LogFormat

Choose how to parse log data. "json" and "plain" are valid. Default value is "plain".

#### Keywords

We can define a filter string to drop useless log.

Filter string use '|' to connect multiple strings. 

`error1|error2`

#### Extractor

#### Sender

