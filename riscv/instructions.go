package riscv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

func ToStringInstrType(instrType int8) string {
	switch instrType {
	case RInstrType:
		return "RInstrType"
	case IInstrType:
		return "IInstrType"
	case SInstrType:
		return "SInstrType"
	case UInstrType:
		return "UInstrType"
	case JInstrType:
		return "JInstrType"
	case IImmInstrType:
		return "IImmInstrType"
	default:
		return fmt.Sprintf("Unknown InstrType (val=%v)", instrType)
	}
}

const (
	RInstrType    int8 = 1
	IInstrType    int8 = 2
	SInstrType    int8 = 3
	UInstrType    int8 = 4
	JInstrType    int8 = 5
	IImmInstrType int8 = 6
)

const (
	STORE    int8 = 0   // 0000011
	MISC_MEM int8 = 15  // 0001111
	OP_IMM   int8 = 19  // 0010011
	AUIPC    int8 = 23  // 0010111
	LOAD     int8 = 35  // 0100011
	OP       int8 = 51  // 0110011
	LUI      int8 = 55  // 0110111
	BRANCH   int8 = 99  // 1100011
	JALR     int8 = 103 // 1100111
	JAL      int8 = 111 // 1101111
	SYSTEM   int8 = 115 // 1110011
)

const (
	FUNC7_RINST_0 int8 = 0
	FUNC7_RINST_1 int8 = 32
)

// IInstr
const (
	FUNC3_ADDI int8 = 0
	FUNC3_SLTI int8 = 1
	FUNC3_ANDI int8 = 2
	FUNC3_ORI  int8 = 3
	FUNC3_XORI int8 = 4
	FUNC3_SLLI int8 = 5
	FUNC3_SRLI int8 = 6
	FUNC3_SRAI int8 = 7
)

type Instruction interface {
	Execute(mem *Memory, regs *Registers) error
	String() string
}

type RInstr struct {
	func7  int8
	rs2    int
	rs1    int
	func3  int8
	rd     int
	opcode int8
}

func (Instr RInstr) String() string {
	return fmt.Sprintf("IInstr{func7=%d, rs2=%d, rs1=%d, func3=%d, rd=%d, opcode=%d}",
		Instr.func7,
		Instr.rs2,
		Instr.rs1,
		Instr.func3,
		Instr.rd,
		Instr.opcode)
}

// RInstr
const (
	FUNC3_ADD  int8 = 0
	FUNC3_SLT  int8 = 2
	FUNC3_SLTU int8 = 3
	FUNC3_AND  int8 = 7
	FUNC3_OR   int8 = 6
	FUNC3_XOR  int8 = 4
	FUNC3_SLL  int8 = 1
	FUNC3_SRL  int8 = 5
	FUNC3_SUB  int8 = 0
	FUNC3_SRA  int8 = 5
)

// RInstr
const (
	FUNC7_ADD  int8 = 0
	FUNC7_SLT  int8 = 0
	FUNC7_SLTU int8 = 0
	FUNC7_AND  int8 = 0
	FUNC7_OR   int8 = 0
	FUNC7_XOR  int8 = 0
	FUNC7_SLL  int8 = 0
	FUNC7_SRL  int8 = 0
	FUNC7_SUB  int8 = 32 // 0100000
	FUNC7_SRA  int8 = 32 // 0100000

)

func (Inst RInstr) Execute(mem *Memory, regs *Registers) error {
	// ADD performs the addition of rs1 and rs2. SUB performs the subtraction of rs2 from rs1. Overflows
	// are ignored and the low XLEN bits of results are written to the destination rd.
	if Inst.opcode == OP {
		if Inst.func7 == FUNC7_ADD && Inst.func3 == FUNC3_ADD {
			// ignore overflow
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] + regs.reg[Inst.rs2]
		} else if Inst.func7 == FUNC7_SUB && Inst.func3 == FUNC3_SUB {
			// ignore overfloat
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] - regs.reg[Inst.rs2]
		} else if Inst.func7 == FUNC7_SLTU && Inst.func3 == FUNC3_SLTU {
			// SLT and SLTU perform signed and unsigned compares respectively, writing 1 to rd if rs1 < rs2, 0 otherwise. Note
			// SLTU rd, x0, rs2 sets rd to 1 if rs2 is not equal to zero, otherwise sets rd to zero (assembler
			// pseudoinstruction SNEZ rd, rs).
			unsigned_rs1 := IntAbs(regs.reg[Inst.rs1])
			unsigned_rs2 := IntAbs(regs.reg[Inst.rs2])
			if unsigned_rs1 < unsigned_rs2 {
				regs.reg[Inst.rd] = 1
			} else {
				regs.reg[Inst.rd] = 0
			}
		} else if Inst.func7 == FUNC7_SLT && Inst.func3 == FUNC3_SLT {
			if regs.reg[Inst.rs1] < regs.reg[Inst.rs2] {
				regs.reg[Inst.rd] = 1
			} else {
				regs.reg[Inst.rd] = 0
			}
		} else if Inst.func7 == FUNC7_AND && Inst.func3 == FUNC3_AND {
			// AND, OR, and XOR perform bitwise logical operations.regs.pc = regs.pc + 1
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] & regs.reg[Inst.rs2]
		} else if Inst.func7 == FUNC7_OR && Inst.func3 == FUNC3_OR {
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] | regs.reg[Inst.rs2]
		} else if Inst.func7 == FUNC7_XOR && Inst.func3 == FUNC3_XOR {
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] ^ regs.reg[Inst.rs2]
		} else if Inst.func7 == FUNC7_SLL && Inst.func3 == FUNC3_SLL {
			// SLL, SRL, and SRA perform logical left, logical right, and arithmetic right shifts on the value in
			// register rs1 by the shift amount held in the lower 5 bits of register rs2.
			// 11111=31
			filter_5_bit := int32(31)
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] << regs.reg[Inst.rs2] & filter_5_bit
		} else if Inst.func7 == FUNC7_SRA && Inst.func3 == FUNC3_SRA {
			// arithmetic shift so keep the sign
			filter_5_bit := int32(31)
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] >> (regs.reg[Inst.rs2] & filter_5_bit)
		} else if Inst.func7 == FUNC7_SRL && Inst.func3 == FUNC3_SRL {
			filter_5_bit := uint32(31)
			// logical one, so we need to convert to unsigned first
			unsigned_rs1 := ReinterpreteAsUnsigned(regs.reg[Inst.rs1])
			unsigned_rs2 := uint32(regs.reg[Inst.rs2]) & filter_5_bit

			unsigned_rd := unsigned_rs1 >> unsigned_rs2
			regs.reg[Inst.rd] = ReinterpreteAsSigned(unsigned_rd)
		}
	}

	return nil
}

type IInstr struct {
	imm    int32
	rs1    int
	func3  int8
	rd     int
	opcode int8
}

func (Instr IInstr) String() string {
	return fmt.Sprintf("IInstr{imm=%d, rs1=%d, func3=%d, rd=%d, opcode=%d}",
		Instr.imm,
		Instr.rs1,
		Instr.func3,
		Instr.rd,
		Instr.opcode)
}

func (Inst IInstr) Execute(mem *Memory, regs *Registers) error {
	switch Inst.opcode {
	case JALR:
		// The indirect jump instruction JALR (jump and link register) uses the I-type encoding. The target
		// address is obtained by adding the sign-extended 12-bit I-immediate to the register rs1, then setting
		// the least-significant bit of the result to zero. The address of the instruction following the jump
		// (pc+4) is written to register rd. Register x0 can be used as the destination if the result is not
		// required.
		regs.reg[Inst.rd] = regs.pc + 1
		regs.pc = regs.reg[Inst.rs1] + Inst.imm
		regs.pc = regs.pc - (regs.pc % 2)
	case OP_IMM:
		return op_imm_execute(Inst, mem, regs)
	default:
		panic("Unknown operator type on IInstr")
	}
	return nil
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
	imm1   uint32
	rs2    uint32
	rs1    uint32
	func3  uint32
	imm0   uint32
	opcode uint32
}

func (Instr SInstr) String() string {
	return fmt.Sprintf("IInstr{imm1=%d, rs2=%d, rs1=%d, func3=%d, imm0=%d, opcode=%d}",
		Instr.imm1,
		Instr.rs2,
		Instr.rs1,
		Instr.func3,
		Instr.imm0,
		Instr.opcode)
}

func (Instr SInstr) Execute(mem *Memory, regs *Registers) error {
	return errors.New("SInstr execute is not implemented")
}

type BInstr struct {
	imm3   uint32
	imm2   uint32
	rs2    uint32
	rs1    uint32
	func3  uint32
	imm1   uint32
	imm0   uint32
	opcode uint32
}

func (Instr BInstr) String() string {
	return fmt.Sprintf(
		"IInstr{imm3=%d, imm2=%d, rs2=%d, rs1=%d, func3=%d, imm1=%d, imm0=%d, opcode=%d}",
		Instr.imm3,
		Instr.imm2,
		Instr.rs2,
		Instr.rs1,
		Instr.func3,
		Instr.imm1,
		Instr.imm0,
		Instr.opcode)
}

func (Instr BInstr) Execute(mem *Memory, regs *Registers) error {
	return errors.New("BInstr execute is not implemented")
}

type UInstr struct {
	imm    int32 // 20 bit offset, I think it should be signed, or you can't jump backwards, but not sure.
	rd     int32
	opcode int8
}

func (Instr UInstr) String() string {
	return fmt.Sprintf(
		"IInstr{imm=%d, rd=%d, opcode=%d",
		Instr.imm,
		Instr.rd,
		Instr.opcode)
}

func (Inst UInstr) Execute(mem *Memory, regs *Registers) error {
	imm1_shifted := Inst.imm << 12 // fill lowest 12 bits with zero
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

type JInstr struct {
	imm3   uint32 // 1 bit
	imm2   uint32 // 10
	imm1   uint32 // 1 bit
	imm0   uint32 // 8 bit immediate
	rd     int
	opcode int8
}

func (Instr JInstr) String() string {
	return fmt.Sprintf(
		"IInstr{(imm3=%d, imm2=%d, imm1=%d, imm0=%d) imm=%d, rd=%d, opcode=%d}",
		Instr.imm3,
		Instr.imm2,
		Instr.imm1,
		Instr.imm0,
		Instr.Imm(),
		Instr.rd,
		Instr.opcode)
}

func (Instr JInstr) Imm() int32 {
	// imm0 = 8 bit
	// imm1 = 1 bits
	// imm2 = 10 bit
	// imm3 = 1 bit

	// imm3::imm0::imm1::imm2
	unsigned := (Instr.imm3 << (1 + 10 + 1 + 8)) +
		(Instr.imm0 << (1 + 10 + 1)) +
		(Instr.imm1 << (1 + 10)) +
		(Instr.imm2 << 1)

	signed := sext(unsigned, 20)

	return ReinterpreteAsSigned(signed)
}

func (Instr JInstr) Execute(mem *Memory, regs *Registers) error {
	if Instr.opcode == JAL {
		// The jump and link (JAL) instruction uses the J-type format, where the J-immediate encodes a
		// signed offset in multiples of 2 bytes. The offset is sign-extended and added to the address of the
		// jump instruction to form the jump target address. Jumps can therefore target a ±1 MiB range.
		// JAL stores the address of the instruction following the jump (pc+4) into register rd. The standard
		// software calling convention uses x1 as the return address register and x5 as an alternate link register.

		// save the next pc to the link register
		regs.reg[Instr.rd] = regs.pc + 1

		// jump to the new location
		regs.pc = regs.pc + Instr.Imm()
	}

	return nil
}

type PInstr struct {
	imm1 int
	imm2 int
	imm3 int
	imm4 int
	rsd  int
}

func CreateADDI(src int, dst int, imm int32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_ADDI, opcode: OP_IMM}
}

func CreateSLLI(src int, dst int, imm int32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SLLI, opcode: OP_IMM}
}

func CreateSLRI(src int, dst int, imm int32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SRLI, opcode: OP_IMM}
}

func CreateSRAI(src int, dst int, imm int32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SRAI, opcode: OP_IMM}
}

func CreateMV(src int, dst int) IInstr {
	return CreateADDI(src, dst, 0)
}

func Nop() IInstr {
	// ADDI x0, x0, 0
	return CreateADDI(0, 0, 0)
}

func CreateLui(imm int32, dst int32) UInstr {
	return UInstr{imm: imm, rd: dst, opcode: LUI}
}

func CreateAUIPC(imm int32, dst int32) UInstr {
	return UInstr{imm: imm, rd: dst, opcode: AUIPC}
}

func CreateADD(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_ADD, func7: FUNC7_ADD, opcode: OP}
}

func CreateSUB(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_SUB, func7: FUNC7_SUB, opcode: OP}
}

func CreateSLTU(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_SLTU, func7: FUNC7_SLTU, opcode: OP}
}

func CreateSLT(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_SLT, func7: FUNC7_SLT, opcode: OP}
}

func CreateAND(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_AND, func7: FUNC7_AND, opcode: OP}
}

func CreateOR(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_OR, func7: FUNC7_OR, opcode: OP}
}

func CreateXOR(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_XOR, func7: FUNC7_XOR, opcode: OP}
}

func CreateSLL(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_SLL, func7: FUNC7_SLL, opcode: OP}
}

func CreateSRA(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_SRA, func7: FUNC7_SRA, opcode: OP}
}

func CreateSRL(rd int, rs1 int, rs2 int) RInstr {
	return RInstr{rd: rd, rs1: rs1, rs2: rs2, func3: FUNC3_SRL, func7: FUNC7_SRL, opcode: OP}
}

func CreateJAL(imm int32, link_reg int) JInstr {
	if imm%2 == 1 {
		panic("imm in jal must be an even number")
	}

	// imm0 = 8 bit
	// imm1 = 1 bits
	// imm2 = 10 bit
	// imm3 = 1 bit

	// imm3::imm0::imm1::imm2::0
	// ..00001010
	word := ReinterpreteAsUnsigned(imm)
	imm0 := bitSliceBetween(word, 12, 19)
	imm1 := bitSliceBetween(word, 11, 11)
	imm2 := bitSliceBetween(word, 1, 10)
	imm3 := bitSliceBetween(word, 20, 20)

	return JInstr{imm0: imm0, imm1: imm1, imm2: imm2, imm3: imm3, rd: link_reg, opcode: JAL}
}

func CreateJ(imm int32) JInstr {
	return CreateJAL(imm, reg_zero)
}

func CreateJALR(imm int32, link_reg int, addr_reg int) IInstr {
	return IInstr{imm: imm, rd: link_reg, opcode: JALR, rs1: addr_reg}
}
