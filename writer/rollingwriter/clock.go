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

package rollingwriter

import (
	"sync"
	"time"
)

type Clocker interface {
	Now() time.Time
}

var (
	_ Clocker = (*WallClocker)(nil)
	_ Clocker = (*FakeClocker)(nil)
)

type WallClocker struct {
}

func (wc WallClocker) Now() time.Time { return time.Now() }

type FakeClocker struct {
	sync.RWMutex
	now time.Time
}

func (f *FakeClocker) SetNow(n time.Time) {
	f.Lock()
	f.now = n
	f.Unlock()
}

func (f *FakeClocker) Now() time.Time {
	f.RLock()
	defer f.RUnlock()
	return f.now
}
