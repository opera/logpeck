#!/bin/bash

function Usage() {
 echo "Usage:"
 echo "  $0 [add|remove|stop|start|update]"
}

cmd=""
case $1 in
	add|remove|stop|start|update)
	 	cmd=$1
		;;
	*)
		Usage; exit 1
		;;
esac

url="http://127.0.0.1:7117/peck_task/$cmd"
config='{
  "Name":"TestLog",
	"LogPath":".test.log",
	"FilterExpr":"mocklog hahaha|mocklog",
	"ESConfig":{
	  "URL":"http://127.0.0.1:9200",
		"Index":"mocklog10",
		"Type":"Mocks"
	}
}'
#echo $url
#echo $config
curl -XPOST $url -d "$config"
