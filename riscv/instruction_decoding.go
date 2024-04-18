package riscv

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func DecodeIInstr(word uint32) IInstr {
	imm := bitSliceBetween(word, 20, 31)
	rs1 := bitSliceBetween(word, 15, 19)
	func3 := bitSliceBetween(word, 12, 14)
	rd := bitSliceBetween(word, 7, 11)
	opcode := bitSliceBetween(word, 0, 6)

	return IInstr{
		imm:    int32(sext(imm, 12)),
		rs1:    int(rs1),
		func3:  int8(func3),
		rd:     int(rd),
		opcode: int8(opcode),
	}
}

func DecodeSInstr(word uint32) Instruction {
	imm1 := bitSliceBetween(word, 25, 31)
	rs2 := bitSliceBetween(word, 20, 24)
	rs1 := bitSliceBetween(word, 15, 19)
	func3 := bitSliceBetween(word, 12, 14)
	imm0 := bitSliceBetween(word, 7, 11)
	opcode := bitSliceBetween(word, 0, 6)

	return SInstr{
		imm1:   imm1,
		rs2:    rs2,
		rs1:    rs1,
		func3:  func3,
		imm0:   imm0,
		opcode: opcode,
	}
}

func (Instr BInstr) imm() uint32 {
	return Instr.imm0 +
		(Instr.imm1 << 1) +
		(Instr.imm2 << (4 + 1)) +
		(Instr.imm3 << (4 + 1 + 6))
}

func DecodeBInstr(word uint32) BInstr {
	imm3 := bitSliceBetween(word, 31, 31) // 1 bit
	imm2 := bitSliceBetween(word, 25, 30) // 6 bit
	rs2 := bitSliceBetween(word, 20, 24)
	rs1 := bitSliceBetween(word, 15, 19)
	func3 := bitSliceBetween(word, 12, 14)
	imm1 := bitSliceBetween(word, 8, 11) // 4 bit
	imm0 := bitSliceBetween(word, 7, 7)  // 1 bit
	opcode := bitSliceBetween(word, 0, 6)

	// imm = 1+6+4+1 = 12 bit

	return BInstr{
		imm3:   imm3,
		imm2:   imm2,
		rs2:    rs2,
		rs1:    rs1,
		func3:  func3,
		imm1:   imm1,
		imm0:   imm0,
		opcode: opcode,
	}
}

func DecodeUInstr(word uint32) UInstr {
	imm := bitSliceBetween(word, 12, 31)
	rd := bitSliceBetween(word, 7, 11)
	opcode := bitSliceBetween(word, 0, 6)

	return UInstr{
		imm:    ReinterpreteAsSigned(sext(imm, 12)),
		rd:     int32(rd),
		opcode: int8(opcode),
	}
}

func DecodeJInstr(word uint32) Instruction {
	imm3 := bitSliceBetween(word, 31, 31) // 1 bit
	imm2 := bitSliceBetween(word, 21, 30) // 10 bit
	imm1 := bitSliceBetween(word, 20, 20) // 1 bit
	imm0 := bitSliceBetween(word, 12, 19) // 8 bit
	rd := bitSliceBetween(word, 7, 11)
	opcode := bitSliceBetween(word, 0, 6)

	return JInstr{
		imm3:   imm3,
		imm2:   imm2,
		imm1:   imm1,
		imm0:   imm0,
		rd:     int(rd),
		opcode: int8(opcode)}
}

func DecodeRInstr(word uint32) Instruction {
	func7 := bitSliceBetween(word, 25, 31)
	rs2 := bitSliceBetween(word, 20, 24)
	rs1 := bitSliceBetween(word, 15, 19)
	func3 := bitSliceBetween(word, 12, 14)
	rd := bitSliceBetween(word, 7, 11)
	opcode := bitSliceBetween(word, 0, 6)

	return RInstr{
		func7:  int8(func7),
		rs2:    int(rs2),
		rs1:    int(rs1),
		func3:  int8(func3),
		rd:     int(rd),
		opcode: int8(opcode),
	}
}

func DecodePInstr(word uint32) Instruction { panic("Decoding of P instruction not implemented.") }

type Decoder struct {
	OpcodeToInstrType map[int8]int8
}

func NewDecoder() *Decoder {
	d := new(Decoder)
	d.OpcodeToInstrType = map[int8]int8{}

	return d
}

func (d *Decoder) RegisterBaseInstructionSet() {
	if d.OpcodeToInstrType == nil {
		panic("invalid decoder no map present")
	}
	d.Register(OP_IMM, IInstrType)
	d.Register(LUI, UInstrType)
	d.Register(AUIPC, UInstrType)
	d.Register(OP, RInstrType)
	d.Register(JAL, JInstrType)
	d.Register(JALR, IInstrType)
	d.Register(BRANCH, IInstrType)
	d.Register(LOAD, IInstrType)
	d.Register(STORE, SInstrType)
	d.Register(SYSTEM, IInstrType)
	// d.Register(MISC_MEM, ??? ) -> TODO: add when implementin git

}

func (d *Decoder) Register(opcode int8, instrType int8) error {
	unexpectedEntry, alreadyPresent := d.OpcodeToInstrType[opcode]
	if alreadyPresent {
		return fmt.Errorf("opcode %v already registered in decoder", ToStringInstrType(unexpectedEntry))
	}
	d.OpcodeToInstrType[opcode] = instrType

	return nil
}

func ByteArrayToWord(wordArray [4]byte) uint32 {
	word := uint32(0)

	buf := bytes.NewBuffer(wordArray[:])
	binary.Read(buf, binary.LittleEndian, &word) // write takes ownership of buf

	return word
}

func (d Decoder) Decode(word uint32) (Instruction, error) {
	opcode := int8(bitSliceBetween(word, 0, 6))
	instrType, isPresent := d.OpcodeToInstrType[opcode]

	if !isPresent {
		return nil, fmt.Errorf("opcode (%v) not registered in the decoder", opcode)
	}

	switch instrType {
	case RInstrType:
		return DecodeRInstr(word), nil
	case IInstrType:
		return DecodeIInstr(word), nil
	case SInstrType:
		return DecodeSInstr(word), nil
	case UInstrType:
		return DecodeUInstr(word), nil
	case JInstrType:
		return DecodeJInstr(word), nil
	}

	return nil, fmt.Errorf("invalid instrtype(%v) ", ToStringInstrType(instrType))
}
