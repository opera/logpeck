# Peck Task Configuration

Logpeck use a json string to define a task. A simplest config is as follows.

```
{
    "Name":"HttpServer",
    "LogPath":"/var/log/http_server/http_server.out.log",
    "ESConfig":{
        "Hosts":["172.10.1.1:9200","172.10.1.2:9200","172.10.1.3:9200"],
        "Index":"http_server-%{+2006.01.02}",
        "Type":"ErrorLog"
    }
}
```

## Required Configuration

There are three required field named "Name", "LogPath", "ESConfig". 
If we only use required configuration, Logpeck will post each line from "LogPath" into ElasticSearch.

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

#### FilterExpr

We can define a filter string to drop useless log.

Filter string use '|' to connect multiple strings. 

`error1|error2`

#### Delimiters

#### Fields

