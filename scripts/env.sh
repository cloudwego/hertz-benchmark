#!/bin/bash

nprocs=$(getconf _NPROCESSORS_ONLN)
if [ $nprocs -lt 4 ]; then
  echo "Your environment should have at least 4 processors"
  exit 1
elif [ $nprocs -gt 20 ]; then
  nprocs=20
fi

# GO
GOEXEC=${GOEXEC:-"go"}
GOROOT=$GOROOT

scpu=$((nprocs > 16 ? 3 : nprocs / 4 - 1)) # max is 3(4 cpus)
taskset_less="taskset -c 0-$scpu"
taskset_more="taskset -c $((scpu + 1))-$((nprocs - 1))"

REPORT=${REPORT:-"$(date +%F-%H-%M)"}
tee_cmd="tee -a output/${REPORT}.log"
