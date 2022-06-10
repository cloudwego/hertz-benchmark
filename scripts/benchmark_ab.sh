#!/bin/bash

. ./scripts/env.sh
# env for collecting data
REPORT_PATH="output/${REPORT}.log"
DATA_PATH="output/${REPORT}.csv"
LATENCY_PATH="output/latency_${REPORT}"

n=3000000
body=(1024 2048 4096 8192 16384 32768 65536)
concurrent=(100)
header=(1024)
repo=("hertz" "fasthttp" "gin" "fasthttp_timeout")
ports=(8001 8002 8003 8004)
serverIP="http://127.0.0.1"

. ./scripts/build_all.sh

# make folder to store all latency data
if [ ! -d ${LATENCY_PATH} ]; then
  mkdir ${LATENCY_PATH}
fi

# generate request in json for ab
for b in ${body[@]}; do
    python ./scripts/ab/generate_request.py ${b}
done

# benchmark
for b in ${body[@]}; do
  for h in ${header[@]}; do
    for c in ${concurrent[@]}; do
      for ((i = 0; i < ${#repo[@]}; i++)); do
        rp=${repo[i]}
        addr="${serverIP}:${ports[i]}"
        # server start
        
        nohup $taskset_less ./output/bin/${rp}_server >>output/log/nohup.log 2>&1 &
        # nohup ./output/bin/${rp}_server >>output/log/nohup.log 2>&1 &
        sleep 1
        echo "server $rp running with $taskset_less"
        # echo "server $rp"

        # run ab
        echo "Benchmark_Config" >> ${REPORT_PATH}
        echo "${rp},${c},${b}" >> ${REPORT_PATH}
        latency_file="${LATENCY_PATH}/${rp}_${c}_${b}.csv"
        $taskset_more ab -e ${latency_file} -d -S -q -n ${n} -c ${c} -k -p ./scripts/ab/request_${b}.json -T application/json ${addr}/ | $tee_cmd

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
python ./scripts/ab/parse_data.py ${REPORT_PATH} ${LATENCY_PATH} ${DATA_PATH}