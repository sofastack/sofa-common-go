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
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func X509(key string, x *x509.Certificate) Field {
	return zap.String(key, x.Subject.String()+" "+x.Issuer.String())
}

func Conn(key string, conn net.Conn) Field {
	return zap.String(key, conn.LocalAddr().String()+"->"+conn.RemoteAddr().String())
}

func Binary(key string, val []byte) Field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Boolp(key string, val *bool) Field { return zap.Boolp(key, val) }

func Skip() Field {
	return zap.Skip()
}

func NilField(key string) Field { return Reflect(key, nil) }

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

// Complex128 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex128 to
// interface{}).
func Complex128(key string, val complex128) Field {
	return zap.Complex128(key, val)
}

func Complex128p(key string, val *complex128) Field { return zap.Complex128p(key, val) }

// Complex64 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex64 to
// interface{}).
func Complex64(key string, val complex64) Field {
	return zap.Complex64(key, val)
}

func Complex64p(key string, val *complex64) Field {
	return zap.Complex64p(key, val)
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshalling is
// necessarily lazy.
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Float64p(key string, val *float64) Field { return zap.Float64p(key, val) }

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshalling is
// necessarily lazy.
func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Float32p(key string, val *float32) Field { return zap.Float32p(key, val) }

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Intp(key string, val *int) Field { return zap.Intp(key, val) }

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int64p(key string, val *int64) Field { return zap.Int64p(key, val) }

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int32p(key string, val *int32) Field { return zap.Int32p(key, val) }

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) Field {
	return zap.Int16(key, val)
}

func Int16p(key string, val *int16) Field { return zap.Int16p(key, val) }

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) Field {
	return zap.Int8(key, val)
}

func Int8p(key string, val *int8) Field { return zap.Int8p(key, val) }

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return zap.String(key, val)
}

func Stringp(key string, val *string) Field { return zap.Stringp(key, val) }

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uintp(key string, val *uint) Field { return zap.Uintp(key, val) }

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uint64p(key string, val *uint64) Field { return zap.Uint64p(key, val) }

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint32p(key string, val *uint32) Field { return zap.Uint32p(key, val) }

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint16p(key string, val *uint16) Field { return zap.Uint16p(key, val) }

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Uint8p(key string, val *uint8) Field { return zap.Uint8p(key, val) }

// Uintptr constructs a field with the given key and value.
func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

func Uintptrp(key string, val *uintptr) Field { return zap.Uintptrp(key, val) }

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

// Namespace creates a named, isolated scope within the logger's context. All
// subsequent fields will be added to the new namespace.
//
// This helps prevent key collisions when injecting loggers into sub-components
// or third-party libraries.
func Namespace(key string) Field {
	return zap.Namespace(key)
}

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Timep(key string, val *time.Time) Field { return zap.Timep(key, val) }

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) Field {
	return zap.Stack(key)
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Durationp(key string, val *time.Duration) Field { return zap.Durationp(key, val) }

// Object constructs a field with the given key and ObjectMarshaler. It
// provides a flexible, but still type-safe and efficient, way to add map- or
// struct-like user-defined types to the logging context. The struct's
// MarshalLogObject method is called lazily.
func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

// Array constructs a field with the given key and ArrayMarshaler. It provides
// a flexible, but still type-safe and efficient, way to add array-like types
// to the logging context. The struct's MarshalLogArray method is called lazily.
func Array(key string, val ArrayMarshaler) Field {
	return zap.Array(key, val)
}

// Bools constructs a field that carries a slicex of bools.
func Bools(key string, bs []bool) Field {
	return zap.Bools(key, bs)
}

// ByteStrings constructs a field that carries a slicex of []byte, each of which
// must be UTF-8 encoded text.
func ByteStrings(key string, bss [][]byte) Field {
	return zap.ByteStrings(key, bss)
}

// Complex128s constructs a field that carries a slicex of complex numbers.
func Complex128s(key string, nums []complex128) Field {
	return zap.Complex128s(key, nums)
}

// Complex64s constructs a field that carries a slicex of complex numbers.
func Complex64s(key string, nums []complex64) Field {
	return zap.Complex64s(key, nums)
}

// Durations constructs a field that carries a slicex of time.Durations.
func Durations(key string, ds []time.Duration) Field {
	return zap.Durations(key, ds)
}

// Float64s constructs a field that carries a slicex of floats.
func Float64s(key string, nums []float64) Field {
	return zap.Float64s(key, nums)
}

// Float32s constructs a field that carries a slicex of floats.
func Float32s(key string, nums []float32) Field {
	return zap.Float32s(key, nums)
}

// Ints constructs a field that carries a slicex of integers.
func Ints(key string, nums []int) Field {
	return zap.Ints(key, nums)
}

// Int64s constructs a field that carries a slicex of integers.
func Int64s(key string, nums []int64) Field {
	return zap.Int64s(key, nums)
}

// Int32s constructs a field that carries a slicex of integers.
func Int32s(key string, nums []int32) Field {
	return zap.Int32s(key, nums)
}

// Int16s constructs a field that carries a slicex of integers.
func Int16s(key string, nums []int16) Field {
	return zap.Int16s(key, nums)
}

// Int8s constructs a field that carries a slicex of integers.
func Int8s(key string, nums []int8) Field {
	return zap.Int8s(key, nums)
}

// Strings constructs a field that carries a slicex of strings.
func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

// Times constructs a field that carries a slicex of time.Times.
func Times(key string, ts []time.Time) Field {
	return zap.Times(key, ts)
}

// Uints constructs a field that carries a slicex of unsigned integers.
func Uints(key string, nums []uint) Field {
	return zap.Uints(key, nums)
}

// Uint64s constructs a field that carries a slicex of unsigned integers.
func Uint64s(key string, nums []uint64) Field {
	return zap.Uint64s(key, nums)
}

// Uint32s constructs a field that carries a slicex of unsigned integers.
func Uint32s(key string, nums []uint32) Field {
	return zap.Uint32s(key, nums)
}

// Uint16s constructs a field that carries a slicex of unsigned integers.
func Uint16s(key string, nums []uint16) Field {
	return zap.Uint16s(key, nums)
}

// Uint8s constructs a field that carries a slicex of unsigned integers.
func Uint8s(key string, nums []uint8) Field {
	return zap.Uint8s(key, nums)
}

// Uintptrs constructs a field that carries a slicex of pointer addresses.
func Uintptrs(key string, us []uintptr) Field {
	return zap.Uintptrs(key, us)
}

// Errors constructs a field that carries a slicex of errors.
func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}

func NamedError(key string, err error) Field {
	return zap.NamedError(key, err)
}

func ErrorString(err error) Field {
	if err != nil {
		return zap.String("error", err.Error())
	}
	return zap.String("error", "nil")
}

func Error(err error) Field {
	return zap.Error(err)
}

type MapStringString map[string]string

func (mss MapStringString) MarshalLogObject(obj zapcore.ObjectEncoder) error {
	for k, v := range mss {
		obj.AddString(k, v)
	}
	return nil
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}
