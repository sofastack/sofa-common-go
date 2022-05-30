<p align="center">
  <b>
    <span style="font-size:larger;">bufio-rw-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/bufio-rw-go"><img src="https://travis-ci.org/detailyang/bufio-rw-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/api/projects/status/fxikraiaswe6x9rc?svg=true"><img src="https://ci.appveyor.com/api/projects/status/hbpj944ankoy9sh5?svg=true" /></a>
   <a href="https://godoc.org/github.com/detailyang/bufio-rw-go">
      <img src="https://godoc.org/github.com/detailyang/bufio-rw-go?status.svg"/>
   </a>
   <br />
   <b>bufio-rw-go is a copy for go's bufio.Reader/Writer but allow caller install hook before read or write</b>
</p>

```bash
go test -v -benchmem -run="^$" github.com/detailyang/bufio-rw-go/ -bench Benchmark
goos: darwin
goarch: amd64
pkg: github.com/detailyang/bufio-rw-go
BenchmarkReaderCopyOptimal-8      	11162181	       102 ns/op	      16 B/op	       1 allocs/op
BenchmarkReaderCopyUnoptimal-8    	 7016584	       172 ns/op	      32 B/op	       2 allocs/op
BenchmarkReaderCopyNoWriteTo-8    	  287985	      3646 ns/op	   32800 B/op	       3 allocs/op
BenchmarkReaderWriteToOptimal-8   	 3374814	       384 ns/op	      16 B/op	       1 allocs/op
BenchmarkWriterCopyOptimal-8      	10851872	       113 ns/op	      16 B/op	       1 allocs/op
BenchmarkWriterCopyUnoptimal-8    	 8423305	       143 ns/op	      32 B/op	       2 allocs/op
BenchmarkWriterCopyNoReadFrom-8   	  337419	      3375 ns/op	   32800 B/op	       3 allocs/op
BenchmarkReaderEmpty-8            	 1587988	       744 ns/op	    4224 B/op	       3 allocs/op
BenchmarkWriterEmpty-8            	 1722464	       670 ns/op	    4096 B/op	       1 allocs/op
BenchmarkWriterFlush-8            	55174977	        19.9 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/bufio-rw-go	14.202s
```
