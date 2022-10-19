# Hertz-Benchmark

调研其他项目的 benchmark 以及 HTTP 使用场景后，我们觉得 Echo 场景能够反映真实的使用场景。本项目提供若干 HTTP 框架在 Echo 场景下的性能记录。

## 使用说明
### 同机压测
执行前请先确认满足环境要求
### Server
```bash
./scripts/benchmark_server.sh
```
### Client
```bash
./scripts/benchmark_client.sh
```
### Profiling
由于默认压测参数会比较迅速完成一次压测，为了获得更长采集时间，可以手动在 ./scripts/benchmark_server.sh 中调整压测参数 n 大小。
#### Profiling Server
不同 server 的 port 映射参见相应脚本，如:
```shell
cat ./scripts/benchmark_pb.sh

# ...
repo=("hertz" "fasthttp" "gin" "fasthttp_timeout")
ports=(8000 8001 8002 8003 8004)
```
获取到对应 server 端口号后，执行：
```shell
go tool pprof localhost:{port}/debug/pprof/{pprof_type}
```

### [wrk Benchmark](https://github.com/wg/wrk)
你也可以用 wrk 作为发压端，参考下面的命令。
```bash
./scripts/benchmark_wrk.sh

# parse data
# ${input_file} locates in /output/$(date +%F-%H-%M).log
# specify one ${output_file}
python ./scripts/wrk/parse_data.py ${input_file} ${output_file} 

# render images
python ./scripts/reports/render_images.py ${output_file}
```

### [ab Benchmark](https://httpd.apache.org/docs/2.4/programs/ab.html)
你也可以用 ab 作为发压端，参考下面的命令。
```bash
./scripts/benchmark_ab.sh

# parse data
# ${input_log_file} locates in /output/$(date +%F-%H-%M).log
# ${input_latency_file} folder locates in /output/latency_$(date +%F-%H-%M)
# specify one ${output_file}
python ./scripts/ab/parse_data.py ${input_log_file} ${input_latency_file} ${output_file} 

# render images (ab does not provide tp999)
python ./scripts/reports/render_images.py ${output_file}
```

## 环境要求
- OS: Linux
    - 默认依赖了命令 taskset, 限定 client 和 server 运行的 CPU; 如在其他系统执行, 请修改脚本。
- CPU: 推荐配置 >=20核, 最低要求 >=4核
    - 压测脚本默认需要 20核 CPU, 具体在脚本的 taskset -c ... 部分, 可以修改或删除。
## 参考数据
  相关说明:

  该压测数据是在调用端有充分机器资源压满服务端的情况下测试，更侧重于关注服务端性能。后续会提供调用端性能数据情况。
### 配置
- CPU: Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz
    - 运行限定 server 4-CPUs，client 16-CPUS
- OS：Debian 5.4.56.bsk.9-amd64 x86_64 GNU/Linux
- Go: 1.16.5
### 数据 (Echo，100 concurrency，1k Header）
  ![Performance](images/performance.png)
