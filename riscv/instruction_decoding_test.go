package riscv

import (
	"testing"
)

func TestDecodeIInstr(t *testing.T) {
	// addi a0, a1, 10
	addi_encoded := uint32(10847507)
	imm := int32(10)
	rs1 := int(11)
	func3 := int8(0)
	rd := 10
	opcode := int8(19)
	expected := IInstr{imm: imm, rs1: rs1, func3: func3, rd: rd, opcode: opcode}
	res := DecodeIInstr(addi_encoded)

	if expected != res {
		t.Logf("\nres=     %s\nexpected=%s", res.String(), expected.String())
		t.Fail()
	}
}

func TestDecodeSInstr(t *testing.T) {
	// sw x1, 3(x2)
	sw_encoded := uint32(1122723)
	imm1 := uint32(0)
	rs2 := uint32(1)
	rs1 := uint32(2)
	func3 := uint32(2)
	imm0 := uint32(3)
	opcode := uint32(35)
	expected := SInstr{
		imm1:   imm1,
		rs2:    rs2,
		rs1:    rs1,
		func3:  func3,
		imm0:   imm0,
		opcode: opcode}
	res := DecodeSInstr(sw_encoded)

	if expected != res {
		t.Logf("\nres=     %s\nexpected=%s", res.String(), expected.String())
		t.Fail()
	}
}

func TestDecodeBInstr(t *testing.T) {
	// beq x1, x2, -16
	imm3 := uint32(1)
	imm2 := uint32(63)
	rs2 := uint32(2)
	rs1 := uint32(1)
	func3 := uint32(0)
	imm1 := uint32(8)
	imm0 := uint32(1)
	opcode := uint32(99)

	expected := BInstr{
		imm3,
		imm2,
		rs2,
		rs1,
		func3,
		imm1,
		imm0,
		opcode}

	if expected.imm() != 4081 {
		t.Fatalf("expected.imm()(val=%d) != 4081", expected.imm())
	}

	bqe_encoded := uint32(4263545059)

	res := DecodeBInstr(bqe_encoded)

	if expected != res {
		t.Fatalf("\nres=     %s\nexpected=%s", res.String(), expected.String())
	}
}

func TestDecodeUInstr(t *testing.T) {
	// lui x1, 10
	expected := UInstr{
		imm:    10,
		rd:     1,
		opcode: 55,
	}

	lui_encoded := uint32(41143)

	res := DecodeUInstr(lui_encoded)

	if expected != res {
		t.Fatalf("\nres=     %s\nexpected=%s", res.String(), expected.String())
	}
}

func TestDecodeJInstr(t *testing.T) {
	// jal x1, -1042430
	expected := JInstr{
		imm3:   1,
		imm2:   1,
		imm1:   1,
		imm0:   1,
		rd:     1,
		opcode: JAL,
	}

	expectedImm := int32(-1042430)
	expectedDerivedImm := expected.Imm()
	if expectedImm != expectedDerivedImm {
		t.Fatalf("\nresImm=     %v\nexpectedImm=%v", expectedDerivedImm, expectedImm)
	}

	// the last 21 bit should like:
	// 100000001100000000010
	// 1::00000001::1::0000000001::0
	// imm3::imm0::imm1::imm2
	//
	// imm0 = 8 bit
	// imm1 = 1 bits
	// imm2 = 10 bit
	// imm3 = 1 bit

	jal_encoded := uint32(2150633711)

	res := DecodeJInstr(jal_encoded)

	if expected != res {
		t.Fatalf("\nres=     %s\nexpected=%s", res.String(), expected.String())
	}
}

func TestDecodeRInstr(t *testing.T) {
	// add x1, x2, x3
	expected := RInstr{
		func7:  0,
		rs2:    3,
		rs1:    2,
		func3:  0,
		rd:     1,
		opcode: OP,
	}

	// 00000000001100010000000010110011
	add_encoded := uint32(3211443)
	res := DecodeRInstr(add_encoded)

	if expected != res {
		t.Fatalf("\nres=%s\nexpected=%s", res.String(), expected.String())
	}
}

func TestByteArrayToWord(t *testing.T) {
	var input = [4]byte{1, 0, 0, 0}
	res := ByteArrayToWord(input)
	expect := uint32(1)
	Assert(t, res, expect)

	input = [4]byte{0, 1, 0, 0}
	res = ByteArrayToWord(input)
	expect = uint32(256)
	Assert(t, res, expect)
}
