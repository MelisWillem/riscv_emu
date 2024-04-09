package riscv

import "testing"

func Assert[K comparable](t *testing.T, value K, expected K) {
	if value != expected {
		t.Logf("%s: value(%v) != expected(%v) \n", t.Name(), value, expected)
		t.Fail()
	}
}
