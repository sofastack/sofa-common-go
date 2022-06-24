package fnv

import "testing"

// nolint
func TestFnv(t *testing.T) {
	HashAdd(HashNew(), "abcd")
}
