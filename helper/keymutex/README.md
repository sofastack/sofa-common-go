<p align="center">
  <b>
    <span style="font-size:larger;">keymutex-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/keymutex-go"><img src="https://travis-ci.org/detailyang/keymutex-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/keymutex-go"><img src="https://ci.appveyor.com/api/projects/status/hbpj944ankoy9sh5?svg=true" /></a>
   <br />
   <b>keymutex-go is a thread-safe mutex for acquiring locks on arbitrary strings which is from k8s/utils but zero allocate.</b>
</p>

```bash
go test -v -benchmem -run="^$" github.com/detailyang/keymutex-go -bench Benchmark
goos: darwin
goarch: amd64
pkg: github.com/detailyang/keymutex-go
BenchmarkKeyMutex
BenchmarkKeyMutex-8   	31639718	        37.5 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/keymutex-go	1.235s
```
