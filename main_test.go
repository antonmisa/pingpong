package main

import (
	"testing"
)

const N int = 1000

var f = func(s string) {
	return
}

func BenchmarkPingpong_lock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pingpong_lock(N, f)
	}
}

func BenchmarkPingpong_ch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pingpong_ch(N, f)
	}
}

func BenchmarkPingpong_atomic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pingpong_atomic(N, f)
	}
}
