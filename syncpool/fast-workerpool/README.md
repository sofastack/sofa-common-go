<p align="center">
  <b>
    <span style="font-size:larger;">fast-workerpool-go</span>
  </b>
  <br />
   <a href="https://travis-ci.org/detailyang/fast-workerpool-go"><img src="https://travis-ci.org/detailyang/fast-workerpool-go.svg?branch=master" /></a>
   <a href="https://ci.appveyor.com/project/detailyang/fast-workerpool-go"><img src="https://ci.appveyor.com/api/projects/status/ux7lf3h9wf8bx8ep?svg=true" /></a>
   <br />
   <b>fast-workerpool-go ports <a href="https://github.com/valyala/fasthttp">fasthttp</a> FIFO worker pool</b>
</p>

````bash
=== RUN   TestWorkerPool
--- PASS: TestWorkerPool (0.00s)
goos: darwin
goarch: amd64
pkg: github.com/detailyang/fast-worker-pool-go
BenchmarkWorkerPool-8        	 2505943	       423 ns/op	       8 B/op	       1 allocs/op
BenchmarkAntsWorkerPool-8    	 2428760	       436 ns/op	       8 B/op	       1 allocs/op
BenchmarkTunnyWorkerPool-8   	  521528	      2161 ns/op	      26 B/op	       2 allocs/op
BenchmarkSlaveWorkerPool-8   	 2424878	      1135 ns/op	       8 B/op	       0 allocs/op
PASS
ok  	github.com/detailyang/fast-worker-pool-go	34.960s
````
