package fnv

// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Inline and byte-free variant of hash/fnv's fnv64a.

const (
    offset64 = 14695981039346656037
    prime64  = 1099511628211
)

// HashNew initializies a new fnv64a hash value.
func HashNew() uint64 {
    return offset64
}

// HashAdd adds a string to a fnv64a hash value, returning the updated hash.
func HashAdd(h uint64, s string) uint64 {
    for i := 0; i < len(s); i++ {
        h ^= uint64(s[i])
        h *= prime64
    }
    return h
}

// HashAddBytes adds bytes to a fnv64a hash value, returning the updated hash.
func HashAddBytes(h uint64, b []byte) uint64 {
    for i := 0; i < len(b); i++ {
        h ^= uint64(b[i])
        h *= prime64
    }
    return h
}

// HashAddByte adds a byte to a fnv64a hash value, returning the updated hash.
func HashAddByte(h uint64, b byte) uint64 {
    h ^= uint64(b)
    h *= prime64
    return h
}
