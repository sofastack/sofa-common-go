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

package testwriter

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/sofastack/sofa-common-go/writer/dsn"
)

var writers sync.Map

type TestWriter struct {
	sync.Mutex
	path    string
	discard bool
	b       []byte
}

func DelAll() {
	writers.Range(func(k, v interface{}) bool {
		writers.Delete(k)
		return true
	})
}

func Get(path string) (*TestWriter, bool) {
	i, ok := writers.Load(path)
	if !ok {
		return nil, ok
	}

	return i.(*TestWriter), true
}

func Del(path string) {
	writers.Delete(path)
}

func New(d *dsn.DSN) (*TestWriter, string, error) {
	dd := d.GetQuery("discard")
	if dd == "" {
		dd = "false"
	}

	discard, err := strconv.ParseBool(dd)
	if err != nil {
		return nil, "", err
	}

	dd = d.GetQuery("trace")
	if dd == "" {
		dd = "false"
	}

	l := d.GetQuery("level")

	trace, err := strconv.ParseBool(dd)
	if err != nil {
		return nil, "", err
	}

	_, ok := writers.Load(d.GetPath())
	if ok {
		return nil, "", fmt.Errorf("duplicate test writer: %s", d.GetPath())
	}

	w := &TestWriter{
		path:    d.GetPath(),
		discard: discard,
	}

	if trace {
		writers.Store(d.GetPath(), w)
	}

	return w, l, nil
}

func (tw *TestWriter) GetPath() string {
	tw.Lock()
	defer tw.Unlock()
	return tw.path
}

func (tw *TestWriter) GetBuffer() []byte {
	tw.Lock()
	defer tw.Unlock()
	return tw.b
}

func (tw *TestWriter) Write(p []byte) (int, error) {
	tw.Lock()
	defer tw.Unlock()
	if tw.discard {
		return len(p), nil
	}
	tw.b = append(tw.b, p...)
	return len(p), nil
}
