// nolint
// Copyright 20xx The Alipay Authors.
//
// @authors[0]: bingwu.ybw(bingwu.ybw@antfin.com|detailyang@gmail.com)
// @authors[1]: robotx(robotx@antfin.com)
//
// *Legal Disclaimer*
// Within this source code, the comments in Chinese shall be the original, governing version. Any comment in other languages are for reference only. In the event of any conflict between the Chinese language version comments and other language version comments, the Chinese language version shall prevail.
// *法律免责声明*
// 关于代码注释部分，中文注释为官方版本，其它语言注释仅做参考。中文注释可能与其它语言注释存在不一致，当中文注释与其它语言注释存在不一致时，请以中文注释为准。
//
//

package examples

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sofastack/sofa-common-go/writer/rollingwriter"
)

func ExampleRollingwriter() {
	file, err := ioutil.TempFile("/tmp", "rollingwriter")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	o := rollingwriter.NewOption()
	o.SetMaxAge(7).SetMaxBackups(7).SetMaxSize(1024) // 7 days and max 1024MB

	w := rollingwriter.New(file.Name(), o)
	n, err := w.Write([]byte("hello"))
	fmt.Println(n)
	// Output: 5
}

func ExampleTimeRollingWriter() {
	file, err := ioutil.TempFile("/tmp", "rollingwriter")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	trw, err := rollingwriter.NewTimeRollingWriter(file.Name(), &rollingwriter.TimeRollingWriterOption{
		TimeFormat:       rollingwriter.DefaultTimeRollingPerHourFormat,
		TimeRollingNamer: rollingwriter.DefaultTimeRollingNamer,
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, err = trw.Write([]byte("hello world")); err != nil {
		log.Fatal(err)
	}
	// Output:
}
