package riscv

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	R    int = 1
	I    int = 2
	S    int = 3
	U    int = 4
	IImm int = 5
)

const (
	FUNC3_ADDI int = 0
	FUNC3_SLTI int = 1
	FUNC3_ANDI int = 2
	FUNC3_ORI  int = 3
	FUNC3_XORI int = 4
	FUNC3_SLLI int = 5
	FUNC3_SRLI int = 6
	FUNC3_SRAI int = 7
)

type Instruction interface {
	Execute(mem *Memory, regs *Registers) error
	Print()
}

type InvalidInstrction struct {
}

func (Inst InvalidInstrction) Execute(mem *Memory, regs *Registers) {
	panic("Trying to execute invalid instruction...")
}

type RInstr struct {
	rs2 int
	rs1 int
	rd  int
}

type IInstr struct {
	imm    int
	rs1    int
	rd     int
	func3  int
	opcode int
}

func (Inst IInstr) Execute(mem *Memory, regs *Registers) error {
	switch Inst.opcode {
	case OP_IMM:
		return op_imm_execute(Inst, mem, regs)
	default:
		panic("Unknown operator type on IInstr")
	}
}

func op_imm_execute(Inst IInstr, _ *Memory, regs *Registers) error {
	switch Inst.func3 {
	case FUNC3_ADDI:
		// ADDI adds the sign-extended 12-bit immediate to register rs1. Arithmetic overflow is ignored and
		// the result is simply the low XLEN bits of the result. ADDI rd, rs1, 0 is used to implement the MV
		// rd, rs1 assembler pseudoinstruction.
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] + int32(Inst.imm)
	case FUNC3_SLTI:
		// SLTI (set less than immediate) places the value 1 in register rd if register rs1 is less than the sign-
		// extended immediate when both are treated as signed numbers, else 0 is written to rd. SLTIU is
		// similar but compares the values as unsigned numbers (i.e., the immediate is first sign-extended to
		if Inst.rs1 < int(Inst.imm) {
			Inst.rd = 1
		} else {
			Inst.rd = 0
		}
	case FUNC3_ANDI:
		// ANDI, ORI, XORI are logical operations that perform bitwise AND, OR, and XOR on register rs1
		// and the sign-extended 12-bit immediate and place the result in rd. Note, XORI rd, rs1, -1 performs
		// a bitwise logical inversion of register rs1 (assembler pseudoinstruction NOT rd, rs).
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] & int32(Inst.imm)
	case FUNC3_ORI:
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] | int32(Inst.imm)
	case FUNC3_XORI:
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] ^ int32(Inst.imm)
	case FUNC3_SLLI:
		// SLLI is a logical left shift (zeros are shifted into the lower bits)
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] << int32(Inst.imm)
	case FUNC3_SRLI:
		// SRLI is a logical right shift (zeros are shifted into the upper bits);
		// A logical shift also shifts the sign bit, we convert to unsigned
		// in there to make sure the shift also shifts the sign bit.

		buf := new(bytes.Buffer)
		shift := uint32(0)

		// convert to unsigned
		binary.Write(buf, binary.LittleEndian, regs.reg[Inst.rd])
		binary.Read(buf, binary.LittleEndian, &shift)

		// shift it
		shift = shift >> uint32(Inst.imm)

		// convert back to signed
		binary.Write(buf, binary.LittleEndian, shift)
		binary.Read(buf, binary.LittleEndian, &regs.reg[Inst.rs1])
	case FUNC3_SRAI:
		// SRAI is an arithmetic right shift (the original sign bit is copied into the vacated upper bits)
		// we don't cap the input as the sign bit should not be shifted here.
		// so the sext makes sense here.
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] >> int32(Inst.imm)
	default:
		return errors.New("invalid op_imm instruction")
	}
	regs.pc = regs.pc + 1
	return nil
}

type SInstr struct {
	imm   int
	rs1   int
	rs2   int
	func3 int
}

type BInstr struct {
	imm1 int
	imm2 int
	imm3 int
	imm4 int
	rs1  int
	rs2  int
}

type UInstr struct {
	imm1   int32 // 20 bit offset, I think it should be signed, or you can't jump backwards, but not sure.
	rd     int32
	opcode int
}

func (Inst UInstr) Execute(mem *Memory, regs *Registers) error {
	imm1_shifted := Inst.imm1 << 12 // fill lowest 12 bits with zero
	switch Inst.opcode {
	case AUIPC:
		// AUIPC (add upper immediate to pc) is used to build pc-relative addresses and uses the U-type
		// format. AUIPC forms a 32-bit offset from the 20-bit U-immediate, filling in the lowest 12 bits with
		// zeros, adds this offset to the address of the AUIPC instruction, then places the result in register
		// rd.
		regs.reg[Inst.rd] = regs.pc + imm1_shifted
	case LUI:
		// LUI (load upper immediate) is used to build 32-bit constants and uses the U-type format. LUI
		// places the U-immediate value in the top 20 bits of the destination register rd, filling in the lowest
		// 12 bits with zeros.
		regs.reg[Inst.rd] = imm1_shifted
		regs.pc = regs.pc + 1
	default:
		panic("Unknown operator type on IInstr")
	}
	regs.pc = regs.pc + 1
	return nil
}

type PInstr struct {
	imm1 int
	imm2 int
	imm3 int
	imm4 int
	rsd  int
}

const (
	STORE  int = 0
	OP_IMM int = 1
	LOAD   int = 2
	AUIPC  int = 23 // 0010111 = 23
	LUI    int = 55 // 0110111 = 55

)

func CreateADDI(src int, dst int, imm int) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_ADDI, opcode: OP_IMM}
}

func CreateSLLI(src int, dst int, imm int) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SLLI, opcode: OP_IMM}
}

func CreateSLRI(src int, dst int, imm int) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SRLI, opcode: OP_IMM}
}

func CreateSRAI(src int, dst int, imm int) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SRAI, opcode: OP_IMM}
}

func CreateMV(src int, dst int) IInstr {
	return CreateADDI(src, dst, 0)
}

func Nop() IInstr {
	// ADDI x0, x0, 0
	return CreateADDI(0, 0, 0)
}

func Ld(base int, width int, dest int, offset int) IInstr {
	return IInstr{imm: offset, rs1: base, func3: width, rd: dest, opcode: LOAD}
}

func CreateLui(imm int32, dst int32) UInstr {
	return UInstr{imm1: imm, rd: dst, opcode: LUI}
}

func CreateAUIPC(imm int32, dst int32) UInstr {
	return UInstr{imm1: imm, rd: dst, opcode: AUIPC}
}
