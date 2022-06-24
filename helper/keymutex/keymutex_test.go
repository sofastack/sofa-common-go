package keymutex

import "testing"

func TestKeyMutex(t *testing.T) {
	km := New(128)
	km.LockKey("abcd")
	km.UnlockKey("abcd")
}

func BenchmarkKeyMutex(b *testing.B) {
	km := New(128)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		km.LockKey("abcd")
		km.UnlockKey("abcd")
	}
}
