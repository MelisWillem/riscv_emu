package riscv

import "fmt"

func bitSlice(word uint32, numBits uint32, begin uint32) uint32 {
	if numBits == 0 {
		return word
	}
	mask := uint32(1)
	// create mask of n-bits
	for i := uint32(0); i < numBits-1; i++ {
		mask <<= 1
		mask |= uint32(1)
	}
	// shift to the right position
	mask = mask << begin
	// get the value
	ouput_shifted := word & mask
	// shift back to zero bit, and return
	return ouput_shifted >> begin
}

func bitSliceBetween(word uint32, from uint32, to uint32) uint32 {
	if from > to || to > 31 {
		panic(fmt.Sprintf("Invalid bit slice from=%d;to=%d", from, to))
	}
	length := to - from + 1 // +1 as the to index should be included in the range
	return bitSlice(word, length, from)
}

func sext(word uint32, signbitLocation uint32) uint32 {
	isNegative := (word & (uint32(1) << signbitLocation)) > 0
	mask := uint32(0)
	if isNegative {
		for i := signbitLocation + 1; i < 32; i++ {
			mask += 1
			mask <<= 1
		}
		mask <<= signbitLocation
	}

	return word | mask
}
