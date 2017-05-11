#!/bin/bash
cd app && python views.py >../log/logpeck-webserver.log 2>../log/logpeck-webserver.err
