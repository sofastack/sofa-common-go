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

package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
	sofadsn "github.com/sofastack/sofa-common-go/writer/dsn"
	"github.com/sofastack/sofa-common-go/writer/sofawriter"
	"go.uber.org/multierr"
)

var (
	globalRegistry Registry
)

func GetRegistry() *Registry {
	return &globalRegistry
}

type Registry struct {
	sync.RWMutex
	m map[string]*SofaLoggerStatus
}

func NewRegistry() *Registry {
	return &Registry{
		m: make(map[string]*SofaLoggerStatus, 16),
	}
}

func (r *Registry) Sync() error {
	r.RLock()
	defer r.RUnlock()

	var errs []error
	for _, k := range r.m {
		if err := k.GetLogger().Sync(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return multierr.Combine(errs...)
}

func (r *Registry) MarshalJSON() ([]byte, error) {
	r.RLock()
	defer r.RUnlock()

	type Status struct {
		Loggers map[string]*SofaLoggerStatus `json:"loggers"`
	}

	var b bytes.Buffer
	if err := jsoniter.NewEncoder(&b).Encode(&Status{
		Loggers: r.m,
	}); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (r *Registry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		r.DoGet(w, req)
		return
	case "POST", "PUT":
		r.DoPOST(w, req)
		return
	case "":
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (r *Registry) DoGet(w http.ResponseWriter, req *http.Request) {
	_ = jsoniter.NewEncoder(w).Encode(r)
}

func (r *Registry) DoPOST(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	level := req.URL.Query().Get("level")

	r.RLock()
	defer r.RUnlock()
	s, ok := r.m[name]
	if !ok {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(fmt.Sprintf("%s logger not found", name)))
		return
	}

	l := ParseLevel(level)
	s.GetLevel().SetLevel(l)

	w.WriteHeader(200)
	_, _ = w.Write([]byte(l.String()))
}

func (r *Registry) GetLogger(name string) (*SofaLogger, bool) {
	r.RLock()
	defer r.RUnlock()
	sl, ok := r.m[name]
	if ok {
		return sl.GetLogger(), true
	}

	return nil, false
}

func (r *Registry) AllocateLogger(name string, dsn string,
	opts ...Option) (*SofaLogger, error) {
	d, err := sofadsn.NewDSN(dsn)
	if err != nil {
		return nil, err
	}

	r.Lock()
	defer r.Unlock()

	if r.m == nil { // initialize before using
		r.m = make(map[string]*SofaLoggerStatus, 16)
	}

	if _, ok := r.m[name]; ok {
		return nil, fmt.Errorf("duplicated logger name: %s", name)
	}

	writer, err := sofawriter.NewFromDSN(d)
	if err != nil {
		return nil, err
	}

	al := NewAtomicLevel(ParseLevel(writer.GetDSN().GetQuery("level")))

	logger, err := New(writer,
		NewCallerConfig().
			SetName(name).
			SetLevel(al).
			AddOptions(opts...))
	if err != nil {
		_ = writer.Close() // free the writer if need
		return nil, err
	}

	r.m[name] = &SofaLoggerStatus{
		name:   name,
		level:  al,
		writer: writer,
		logger: logger,
	}

	return logger, nil
}

type SofaLoggerStatus struct {
	name   string
	logger *SofaLogger
	level  *AtomicLevel
	writer *sofawriter.Writer
}

func (s *SofaLoggerStatus) GetName() string { return s.name }

func (s *SofaLoggerStatus) GetLogger() *SofaLogger { return s.logger }

func (s *SofaLoggerStatus) GetLevel() *AtomicLevel { return s.level }

func (s *SofaLoggerStatus) GetWriter() *sofawriter.Writer { return s.writer }

func (s *SofaLoggerStatus) MarshalJSON() ([]byte, error) {
	type Status struct {
		Name  string `json:"name"`
		Level string `json:"level"`
		DSN   string `json:"dsn"`
	}

	ms := &Status{}
	ms.Name = s.name
	if mt, err := s.level.MarshalText(); err == nil {
		ms.Level = string(mt)
	}
	ms.DSN = s.writer.GetDSN().String()

	var b bytes.Buffer
	if err := jsoniter.NewEncoder(&b).Encode(ms); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
