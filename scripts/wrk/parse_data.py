import sys
import numpy as np

CONFIG_KEY = "Benchmark_Config"
RESULT_KEY = "Benchmark_Result"

RESULT_LINE_NUM = 5
# RESULT_QPS_KEY = "QPS"
# RESULT_TP99_KEY = "TP99"
# RESULT_TP999_KEY = "TP99.9"
# RESULT_TOATL_REQUEST_KEY = "request"
# RESULT_FAILED_KEY = "failed"

def is_config(line):
    return CONFIG_KEY in line

def is_result(line):
    return RESULT_KEY in line

def parse_config(line):
    l = line.rstrip("\n").split(",")
    return l

def parse_result(lines):
    result = []
    for l in lines:
        result += [l.split(",")[1].rstrip("\n")]
    # 'total request' and 'failed' are not required here
    return result[:-2]


def parse(input, output):
    ret = []
    with open(input, 'r') as f:
        lines = f.readlines()
        config, result = [], []
        for i in range(len(lines)):
            l = lines[i]
            if is_config(l):
                config = parse_config(lines[i+1])
                i += 2
            if is_result(l):
                result = parse_result(lines[i+1:i+RESULT_LINE_NUM+1])
                ret += [','.join(config+result)]
                i += RESULT_LINE_NUM + 1

    with open(output, 'w') as f:
        for l in ret:
            f.write(l)
            f.write('\n')

if __name__ == '__main__':
    if len(sys.argv) > 2:
        input = sys.argv[1]
        output = sys.argv[2]
        parse(input, output)
    else:
        print("Please provide input and output")