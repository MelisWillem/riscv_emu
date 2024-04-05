package riscv

import "testing"

func TestReinterpreteToUnsigned(t *testing.T) {
	// negative number will be very different
	input_negatives := []int32{3, -1, -2}
	// positive numbers should stay the same
	// so 3 -> 3
	// two compliment is inverse + 1
	// so -1 = 111...111
	// -> so negative is 2**32 - 1 = 4294967295
	// so -2 = 111...110
	// -> so negative is 2**23 - 2 = 4294967294
	expected := []uint32{3, 4294967295, 4294967294}

	for i, input := range input_negatives {
		res := ReinterpreteAsUnsigned(input)
		if res != expected[i] {
			t.Logf("res(%d) != expected(%d) ", res, expected[i])
			t.Fail()
		}
	}

}

func TestRoundTripSignedUnsignedReinterprete(t *testing.T) {
	// The smallest number in int32
	input := []int32{1, -1, 2, -2, 2147483647, -2147483648}
	for _, input := range input {
		res_unsigned := ReinterpreteAsUnsigned(input)
		res := ReinterpreteAsSigned(res_unsigned)
		if input != res {
			t.Logf("input(%d) != res(%d) with res_unsigned=%d", input, res, res_unsigned)
			t.Fail()
		}
	}
}

func TestAbs(t *testing.T) {
	// The smallest number in int32
	input := []int32{1, -1, 2, -2, 2147483647, -2147483648}
	// what to do with the most negative number?
	expected := []int32{1, 1, 2, 2, 2147483647, 2147483647}
	for i, input := range input {
		res := IntAbs(input)

		if expected[i] != res {
			t.Logf("expected(%d) != res(%d)", expected[i], res)
			t.Fail()
		}
	}
}
