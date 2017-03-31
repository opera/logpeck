#!/bin/bash

url="http://127.0.0.1:7117/peck_task/add"
config='{"Name":"TestLog","LogPath":".test.log","FilterExpr":"mocklog","ESConfig":{"URL":"http://127.0.0.1:9200","Index":"mocklog","Type":"Timestamp"}}'

curl -XPOST $url -d $config
