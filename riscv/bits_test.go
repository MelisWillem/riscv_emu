package riscv

import (
	"testing"
)

func TestSextNegative(t *testing.T) {
	test_val := uint32(8)              // 00..0001000
	expected_val := uint32(4294967288) // 11..1111000

	res := sext(test_val, 3)

	Assert(t, res, expected_val)
}

func TestSextPositive(t *testing.T) {
	test_val := uint32(8) // 00..0001000

	// we we pass a zero as sign bit, it should do nothing
	res1 := sext(test_val, 4) // after the bit
	res2 := sext(test_val, 2) // before the bit

	Assert(t, test_val, res1)
	Assert(t, test_val, res2)
}

func TestBitSlice(t *testing.T) {
	input := uint32(32 + 16 + 8) // ...00111000

	res_slice_all_1 := bitSliceBetween(input, 3, 5)
	expected_all_1 := uint32(1 + 2 + 4) // ...000111
	Assert(t, res_slice_all_1, expected_all_1)

	res_slice_1_zero_before := bitSliceBetween(input, 2, 5)
	expected_slice_1_zero_before := uint32(2 + 4 + 8) // ...001110
	Assert(t, res_slice_1_zero_before, expected_slice_1_zero_before)
}
