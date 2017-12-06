# Logpeck RESTful API

1. Add a new task first. (more task configuration, see [here](task_config.md))

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

2. Start task.

```
curl -XPOST http://127.0.0.1:7117/peck_task/start -d {
  	"Name":"SystemLog"
}
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

6. List task stats

```
curl -XPOST http://127.0.0.1:7117/peck_task/liststats
```
