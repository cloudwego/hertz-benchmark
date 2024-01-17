#!/bin/bash

. ./scripts/env.sh
# env for collecting data
REPORT_PATH="output/${REPORT}.log"
DATA_PATH="output/${REPORT}.csv"

t=30
body=(1024 2048 4096 8192 16384 32768 65536)
concurrent=(100)
header=(1024)
repo=("hertz" "fasthttp" "gin" "fiber" "hex")
ports=(8001 8002 8003 8004 8005)
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

        # run wrk
        echo "Benchmark_Config" >> ${REPORT_PATH}
        echo "${rp},${c},${b}" >> ${REPORT_PATH}
        $taskset_more wrk -d${t}s -s ./scripts/wrk/benchmark.lua -c${c} -t${c} ${addr} -- ${b} | $tee_cmd

        # stop server
        pid=$(ps -ef | grep ${rp}_server | grep -v grep | awk '{print $2}')
        disown $pid
        kill -9 $pid
        sleep 1
      done
    done
  done
done

# parse data and generate output.csv
python ./scripts/wrk/parse_data.py ${REPORT_PATH} ${DATA_PATH}