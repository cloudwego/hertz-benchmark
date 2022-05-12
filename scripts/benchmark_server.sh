#!/bin/bash

. ./scripts/env.sh

n=3000000
body=(1024 2048 4096 8192 16384 32768 65536)
concurrent=(100)
header=(1024)
repo=("hertz" "fasthttp" "gin" "fasthttp_timeout")
ports=(8001 8002 8003 8004)
serverIP="http://127.0.0.1"

. ./scripts/build_all.sh

# benchmark
for b in ${body[@]}; do
  for h in ${header[@]}; do
    for c in ${concurrent[@]}; do
      for ((i = 0; i < ${#repo[@]}; i++)); do
        rp=${repo[i]}
        addr="${serverIP}:${ports[i]}"
        # server start
        nohup $taskset_less ./output/bin/${rp}_server >>output/log/nohup.log 2>&1 &
        sleep 1
        echo "server $rp running with $taskset_less"

        # run client
        echo "client $rp running with fasthttp_client"
        $taskset_more ./output/bin/fasthttp_client -addr="$addr" -b=$b -h=$h -c=$c -n=$n -s=${rp}_server | $tee_cmd

        # stop server
        pid=$(ps -ef | grep ${rp}_server | grep -v grep | awk '{print $2}')
        disown $pid
        kill -9 $pid
        sleep 1
      done
    done
  done
done
