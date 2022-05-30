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
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTimeRollingWriter(t *testing.T) {
	tmpdir, err := ioutil.TempDir("./testdata/testlog", "")
	require.Nil(t, err)
	defer os.RemoveAll(tmpdir)

	tmpfile, err := ioutil.TempFile(tmpdir, "*.test.log")
	require.Nil(t, err)

	fc := &FakeClocker{}

	trw, err := NewTimeRollingWriter(
		tmpfile.Name(),
		&TimeRollingWriterOption{
			TimeFormat: DefaultTimeRollingPerSecondFormat,
			Clocker:    fc,
		},
	)
	require.Nil(t, err)

	now := time.Now()
	count := 10
	for i := 0; i < count; i++ {
		fc.SetNow(now.Add(time.Duration(i) * time.Second))
		// nolint
		n, ferr := trw.Write([]byte("hello"))
		require.Nil(t, ferr)
		_ = n
	}

	files, err := ioutil.ReadDir(tmpdir)
	require.Nil(t, err)
	require.Equal(t, count, len(files))
}
