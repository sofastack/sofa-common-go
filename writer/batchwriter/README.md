<p align="center">
  <b>
    <span style="font-size:larger;">batchwriter-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/batchwriter-go"><img src="https://travis-ci.org/detailyang/batchwriter-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/batchwriter-go"><img src="https://ci.appveyor.com/api/projects/status/hf9iekr7d9vdq9x4?svg=true" /></a>
   <a href="https://godoc.org/github.com/detailyang/batchwriter-go">
      <img src="https://godoc.org/github.com/detailyang/batchwriter-go?status.svg"/>
   </a>
   <br />
   <b>batchwriter-go implements the io.Writer which batch writes (maybe writev if it's net.Conn) to writer by channel</b>
</p>

```bash
go test -v -benchmem -run="^$" github.com/detailyang/batchwriter-go -bench "^Benchmark"
goos: darwin
goarch: amd64
pkg: github.com/detailyang/batchwriter-go
BenchmarkBatchWrite-8           	 2929246	       369 ns/op	     20997 failure	   2908249 success	     139 B/op	       1 allocs/op
BenchmarkBatchWriteParallel-8   	45410478	        35.0 ns/op	  45330196 failure	     80282 success	       0 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/batchwriter-go	3.131s
```
