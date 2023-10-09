# Hertz-Binding-Benchmark

## Hertz Binding Refactor Benchmark
Reference framework
* [gin](https://github.com/gin-gonic/gin)
* [go-tagexpr](https://github.com/hertz-contrib/binding/tree/main/go_tagexpr)
* [fiber v3](https://github.com/gofiber/fiber/pull/2006) (not yet released, the binding refactor of hertz references its design)

## How to work
| scenario                      | command |
|-------------------------------| ----  |
| less query & less field       | `go test -test.bench="NormalQuery"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |
| more query & less field       | `go test -test.bench="BigQuerySmallField"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |
| more query & more field       | `go test -test.bench="BigQueryBigField"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |
| less query & less slice field | `go test -test.bench="SmallSlice1"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |
| more query & less slice field | `go test -test.bench="BigQuerySmallSlice"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |
| more query & more slice field | `go test -test.bench="BigQueryBigSlice"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |
| less query & more slice field | `go test -test.bench="SmallQueryBigSlice"  -test.benchmem --benchtime=5s   -run="bind_test.go"` |

Note:
* `fiber v3` is not yet released, so the benchmark for it can not be executed. 