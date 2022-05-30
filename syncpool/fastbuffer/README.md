<p align="center">
  <b>
    <span style="font-size:larger;">fastbuffer-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/fastbuffer-go"><img src="https://travis-ci.org/detailyang/fastbuffer-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/fastbuffer-go"><img src="https://ci.appveyor.com/api/projects/status/hbpj944ankoy9sh5?svg=true" /></a>
   <a href="https://godoc.org/github.com/detailyang/fastbuffer-go">
      <img src="https://godoc.org/github.com/detailyang/fastbuffer-go?status.svg"/>
   </a>
   <br />
   <b>fastbuffer-go is type alias to [][]byte but holds the sync.Mutex and do not free []byte to os which can reduce memory allocation.</b>
</p>

````bash
go test -v -benchmem -run="^$" github.com/detailyang/fastbuffer-go -bench "^Benchmark"
goos: darwin
goarch: amd64
pkg: github.com/detailyang/fastbuffer-go
BenchmarkBatchFastBuffer-8   	 3265705	       357 ns/op	  41.96 MB/s	       0 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/fastbuffer-go	1.556s
````
