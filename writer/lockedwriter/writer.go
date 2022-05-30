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

package lockedwriter

import (
	"io"
	"sync"
)

type LockedWriter struct {
	sync.Mutex
	w io.Writer
}

func New(w io.Writer) *LockedWriter {
	return &LockedWriter{w: w}
}

func (lw *LockedWriter) Write(p []byte) (int, error) {
	lw.Lock()
	n, err := lw.w.Write(p)
	lw.Unlock()
	return n, err
}

func (lw *LockedWriter) Close() error {
	if rw, ok := lw.w.(io.Closer); ok {
		return rw.Close()
	}
	return nil
}
