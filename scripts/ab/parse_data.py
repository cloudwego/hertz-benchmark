from keyword import iskeyword
from operator import is_
import sys
import numpy as np

CONFIG_KEY = "Benchmark_Config"
RESULT_QPS_KEY = "Requests per second:"
# LATENCY_FILE_PATH = "./output/latency/"

def is_keyword(line, keyword):
    return keyword in line

def parse_config(line):
    l = line.rstrip("\n").split(",")
    return l

def parse_qps(line):
    suffix = ' [#/sec] (mean)\n'
    return line.lstrip(RESULT_QPS_KEY).lstrip(' ').rstrip(suffix)
    # Requests per second:    31833.62 [#/sec] (mean)

def read_latency_file(path):
    with open(path, 'r') as f:
        lines = f.readlines()
        # tp99 is the second last row of this file
        tp99 = lines[-2].rstrip('\n').split(',')[1]
        return tp99

def parse(log_input, latency_input, output):
    ret = []
    with open(log_input, 'r') as f:
        lines = f.readlines()
        i = 0
        while i < len(lines):
            config = []
            qps, tp99, tp999 = "", "", "0" # no tp999 in log
            for j in range(i, len(lines)):
                l = lines[j]
                # benchmark config
                if is_keyword(l, CONFIG_KEY):
                    config = parse_config(lines[j+1])
                    continue
                # result
                if is_keyword(l, RESULT_QPS_KEY):
                    # make sure the config is complete
                    if len(config) != 3:
                        break
                    # qps
                    qps = parse_qps(l)
                    # tp99, read from file: ${name}_${concurrency}_${body_size}.json
                    file_name = "_".join(config)+".csv"
                    tp99 = read_latency_file(latency_input_path+"/"+file_name)
                    # break after parsing the result
                    break
            if len(config) == 3 and qps != "" and tp99 != "":
                # concat the data
                ret += [','.join(config+[qps, tp99, tp999])]
            i = j + 1

    with open(output, 'w') as f:
        for l in ret:
            f.write(l)
            f.write('\n')

if __name__ == '__main__':
    if len(sys.argv) > 3:
        log_input = sys.argv[1]
        latency_input_path = sys.argv[2]
        output = sys.argv[3]
        parse(log_input, latency_input_path, output)
    else:
        print("Please provide input and output")