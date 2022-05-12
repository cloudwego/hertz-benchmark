#!/bin/bash

ps -ef | grep server | grep -v grep
pid=$(ps -ef | grep server | grep -v grep | awk '{print $2}')
kill -9 $pid
