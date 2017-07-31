# Logpeck

Logpeck is an interactive log collector.

## Features
 * Support plain text and json format log.
 * Full control with HTTP API.
 * Add/Remove/Start/Stop collection task freely.
 * Cooperate with ElasticSearch/Kibana deeply.
 * Collection speed control.
 * Get collection status freely.
 * Web show/control conveniently([logpeck-web](https://github.com/opera/logpeck-web)).
 
## Build & Start

`go build cmd/logpeckd/logpeckd.go`

`./logpeckd -config logpeckd.conf`
 
## Scenarios
 * Interactive analysis/debug
 * Visualized presentation
