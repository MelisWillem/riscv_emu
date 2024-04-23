package riscv

import (
	"bytes"
	"encoding/binary"
	"math"
)

func ReinterpreteAsUnsigned(in int32) uint32 {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, in)

	out := uint32(0)
	binary.Read(buf, binary.LittleEndian, &out)

	return out
}

func ReinterpreteAsSigned(in uint32) int32 {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, in)

	out := int32(0)
	binary.Read(buf, binary.LittleEndian, &out)

	return out
}

func IntAbs(in int32) int32 {
	if in == -2147483648 {
		// special case, this number will overflow if abs
		// for now we will just round to the nearest number.
		// Maybe I should return an error and let the called decide on this.
		return 2147483647 // closest thing to the answer we have
	}

	if math.Signbit(float64(in)) {
		return in * -1
	}
	return in
}

func pow(num uint32, power uint32) uint32 {
	return uint32(math.Round(math.Pow(float64(num), float64(power))))
}
