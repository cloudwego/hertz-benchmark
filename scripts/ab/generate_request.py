import json
import sys

def generate_json(size):
    req = {
        "body": ''.join('0' for _ in range(size))
    }
    with open('./scripts/ab/request_{}.json'.format(size), 'w') as json_file:
        json.dump(req, json_file, ensure_ascii = False)


if __name__ == "__main__":
    if len(sys.argv) > 1:
        body_size = int(sys.argv[1])
        generate_json(body_size)

    # body_sizes = [1024, 2048, 4096, 8192, 16384, 32768, 65536]
    # for s in body_sizes:
        # generate_json(s)