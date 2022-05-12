#!/bin/bash

. ./scripts/env.sh
repo=("hertz" "fasthttp" "net")
ports=8001

. ./scripts/build_all.sh
n=3000000

body=(1024)
concurrent=(100)
serverIP="http://127.0.0.1"

# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for ((i = 0; i < ${#repo[@]}; i++)); do
      rp=${repo[i]}
      addr="${serverIP}:${ports}"
      echo $taskset_more ./output/bin/hertz_server
      # server start
      nohup $taskset_more ./output/bin/hertz_server >>output/log/nohup.log 2>&1 &
      sleep 1
      echo "server hertz running with $taskset_more"
      # run client
      echo "client $rp running with ${rp}_client"
      $taskset_less ./output/bin/${rp}_client -addr="$addr" -b=$b -c=$c -n=$n -s=hertz_server | $tee_cmd

      # stop server
      pid=$(ps -ef | grep hertz_server | grep -v grep | awk '{print $2}')
      disown $pid
      kill -9 $pid
      sleep 1
    done
  done
done
