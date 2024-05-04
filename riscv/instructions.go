package riscv

import (
	"fmt"
	"log"
)

func unknowOpcodeError(opcode int8, instrType int8) error {
	return fmt.Errorf("uknown opcode=%v in InstrType=%v", opcode, ToStringInstrType(instrType))
}

func unknownFunc3Error(func3 int8, opcode int8, instrType int8) error {
	return fmt.Errorf("uknown func3=%v in InstrType=%v with opcode=%v", func3, ToStringInstrType(instrType), opcode)
}

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
	Execute(mem Memory, regs *Registers) error
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

func (Inst RInstr) Execute(mem Memory, regs *Registers) error {
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
			unsigned_rs1 := IntAbs(ReinterpreteAsSigned(regs.reg[Inst.rs1]))
			unsigned_rs2 := IntAbs(ReinterpreteAsSigned(regs.reg[Inst.rs2]))
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
			filter_5_bit := uint32(31)
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] << regs.reg[Inst.rs2] & filter_5_bit
		} else if Inst.func7 == FUNC7_SRA && Inst.func3 == FUNC3_SRA {
			// arithmetic shift so keep the sign
			filter_5_bit := uint32(31)
			regs.reg[Inst.rd] = regs.reg[Inst.rs1] >> (regs.reg[Inst.rs2] & filter_5_bit)
		} else if Inst.func7 == FUNC7_SRL && Inst.func3 == FUNC3_SRL {
			filter_5_bit := uint32(31)
			// logical one, so we need to convert to unsigned first
			unsigned_rs1 := regs.reg[Inst.rs1]
			unsigned_rs2 := uint32(regs.reg[Inst.rs2]) & filter_5_bit

			unsigned_rd := unsigned_rs1 >> unsigned_rs2
			regs.reg[Inst.rd] = unsigned_rd
		}
	} else {
		return unknowOpcodeError(Inst.opcode, RInstrType)
	}

	return nil
}

type IInstr struct {
	imm    uint32
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

const (
	FUNC3_LB  int8 = 0
	FUNC3_LH  int8 = 1
	FUNC3_LW  int8 = 2
	FUNC3_LBU int8 = 4
	FUNC3_LHU int8 = 5
)

func (Inst IInstr) Execute(mem Memory, regs *Registers) error {
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
	case LOAD:
		// Loads are encoded in the I-type format and stores are S-type. The effective address is obtained by
		// adding register rs1 to the sign-extended 12-bit offset. Loads copy a value from memory to register rd.
		// Stores copy the value in register rs2 to memory.
		addr := regs.reg[Inst.rs1] + sext(uint32(Inst.imm), 11)

		var err error
		var val uint32
		switch Inst.func3 {
		case FUNC3_LW:
			// The LW instruction loads a 32-bit value from memory into rd.
			val, err = mem.Load(uint32(addr), 4)
		case FUNC3_LH:
			// LH loads a 16-bit value from memory, then sign-extends to 32-bits before storing in rd.
			val, err = mem.Load(uint32(addr), 2)
			val = sext(val, 15)
		case FUNC3_LHU:
			// LHU loads a 16-bit value from memory but then zero extends to 32-bits before storing in rd.
			val, err = mem.Load(uint32(addr), 2)
		case FUNC3_LB:
			// LB and LBU are defined analogously for 8-bit values.
			val, err = mem.Load(uint32(addr), 1)
			val = sext(val, 7)
		case FUNC3_LBU:
			val, err = mem.Load(uint32(addr), 1)
		}

		if err != nil {
			return err
		}
		regs.reg[Inst.rd] = val
	default:
		return unknowOpcodeError(Inst.opcode, IInstrType)
	}
	return nil
}

func op_imm_execute(Inst IInstr, _ Memory, regs *Registers) error {
	switch Inst.func3 {
	case FUNC3_ADDI:
		// ADDI adds the sign-extended 12-bit immediate to register rs1. Arithmetic overflow is ignored and
		// the result is simply the low XLEN bits of the result. ADDI rd, rs1, 0 is used to implement the MV
		// rd, rs1 assembler pseudoinstruction.
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] + Inst.imm
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
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] & Inst.imm
	case FUNC3_ORI:
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] | Inst.imm
	case FUNC3_XORI:
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] ^ Inst.imm
	case FUNC3_SLLI:
		// SLLI is a logical left shift (zeros are shifted into the lower bits)
		imm_static := bitSliceBetween(Inst.imm, 5, 11)
		if imm_static != 0 {
			return fmt.Errorf("invalid SLLI instruction, the imm[11:5] should be equal to 0 but is %d", imm_static)
		}
		imm_shamt := bitSliceBetween(Inst.imm, 0, 4)
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] << int32(imm_shamt)
	case FUNC3_SRLI:
		// SRLI is a logical right shift (zeros are shifted into the upper bits);
		// A logical shift also shifts the sign bit, we convert to unsigned
		// in there to make sure the shift also shifts the sign bit.
		imm_static := bitSliceBetween(Inst.imm, 5, 11)
		if imm_static != 0 {
			return fmt.Errorf("invalid SRLI instruction, the imm[11:5] should be equal to 0 but is %d", imm_static)
		}
		imm_shamt := bitSliceBetween(Inst.imm, 0, 4)
		regs.reg[Inst.rs1] = regs.reg[Inst.rd] >> imm_shamt
	case FUNC3_SRAI:
		// SRAI is an arithmetic right shift (the original sign bit is copied into the vacated upper bits)
		// We don't cap the input as the sign bit should not be shifted here.
		// so the sext makes sense here.
		imm_static := bitSliceBetween(Inst.imm, 5, 11)
		if imm_static != 32 {
			return fmt.Errorf("invalid SRLI instruction, the imm[11:5] should be equal to 32 but is %d", imm_static)
		}
		imm_shamt := bitSliceBetween(Inst.imm, 0, 4)
		val := ReinterpreteAsSigned(regs.reg[Inst.rd]) >> int32(imm_shamt)
		regs.reg[Inst.rs1] = ReinterpreteAsUnsigned(val)
	default:
		return fmt.Errorf("invalid func3(val=%v) value on op_imm instruction", Inst.func3)
	}
	regs.pc = regs.pc + 1
	return nil
}

type SInstr struct {
	imm1   uint32
	rs2    int
	rs1    int
	func3  int8
	imm0   uint32
	opcode int8
}

func (Instr SInstr) imm() uint32 {
	// imm0 -> offset[0:4]
	// imm1 -> offset[5:11]
	return Instr.imm0 + (Instr.imm1 << 5)
}

func (Instr SInstr) String() string {
	return fmt.Sprintf("IInstr{imm1=%d, rs2=%d, rs1=%d, func3=%d, imm0=%d, opcode=%d} with .imm()=%d",
		Instr.imm1,
		Instr.rs2,
		Instr.rs1,
		Instr.func3,
		Instr.imm0,
		Instr.opcode,
		Instr.imm())
}

const (
	FUNC3_SB int8 = 0
	FUNC3_SH int8 = 1
	FUNC3_SW int8 = 2
)

func (Instr SInstr) Execute(mem Memory, regs *Registers) error {
	switch Instr.opcode {
	case STORE:
		// Loads are encoded in the I-type format and stores are S-type. The effective address is obtained by
		// adding register rs1 to the sign-extended 12-bit offset. Loads copy a value from memory to register rd.
		// Stores copy the value in register rs2 to memory.
		addr := regs.reg[Instr.rs1] + sext(Instr.imm(), 11)
		var val uint32
		var err error
		switch Instr.func3 {
		// The SW, SH, and SB instructions store 32-bit, 16-bit, and 8-bit values from the low bits of register
		// rs2 to memory
		case FUNC3_SB:
			val, err = mem.Load(uint32(addr), 1)
		case FUNC3_SH:
			val, err = mem.Load(uint32(addr), 2)
		case FUNC3_SW:
			val, err = mem.Load(uint32(addr), 4)
		default:

			return unknownFunc3Error(Instr.func3, Instr.opcode, SInstrType)
		}

		if err != nil {
			return err
		}
		regs.reg[Instr.rs2] = val
	}
	return unknowOpcodeError(Instr.opcode, SInstrType)
}

type BInstr struct {
	imm3   uint32
	imm2   uint32
	rs2    uint32
	rs1    uint32
	func3  uint32
	imm1   uint32
	imm0   uint32
	opcode int8
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

const (
	FUNC3_BEQ  uint32 = 0
	FUNC3_BNE  uint32 = 1
	FUNC3_BLT  uint32 = 4
	FUNC3_BLTU uint32 = 6
	FUNC3_BGE  uint32 = 5
	FUNC3_BGEU uint32 = 7
)

func (Instr BInstr) Execute(mem Memory, regs *Registers) error {
	offset := sext(Instr.imm(), 11)
	// rs1 and rs2 are stored as uint32, but in this case they represent
	// a signed number. By shifting we can remove the sign bit before comparing
	// and so get an unsigned expressed
	rs1_unsigned := (Instr.rs1 << 1) >> 1
	rs2_unsigned := (Instr.rs1 << 1) >> 1

	// If the comparison is gte/tle then we need to take into
	// the sign bit:
	rs1_signed := ReinterpreteAsSigned(Instr.rs1)
	rs2_signed := ReinterpreteAsSigned(Instr.rs2)

	switch Instr.func3 {
	// BEQ and BNE take the branch if registers rs1 and rs2
	// are equal or unequal respectively.
	case FUNC3_BEQ:
		if regs.reg[Instr.rs1] == regs.reg[Instr.rs2] {
			regs.pc = regs.pc + offset
		}
	case FUNC3_BNE:
		if regs.reg[Instr.rs1] != regs.reg[Instr.rs2] {
			regs.pc = regs.pc + offset
		}
	// BLT and BLTU take the branch if rs1 is less than rs2, using
	// signed and unsigned comparison respectively.
	case FUNC3_BLT:
		if regs.reg[rs1_signed] < regs.reg[rs2_signed] {
			regs.pc = regs.pc + offset
		}
	case FUNC3_BLTU:
		if regs.reg[rs1_unsigned] < regs.reg[rs2_unsigned] {
			regs.pc = regs.pc + offset
		}
	// BGE and BGEU take the branch if rs1 is greater
	// than or equal to rs2, using signed and unsigned comparison respectively.
	case FUNC3_BGE:
		if regs.reg[rs1_signed] >= regs.reg[rs2_signed] {
			regs.pc = regs.pc + offset
		}
	case FUNC3_BGEU:
		if regs.reg[rs1_unsigned] >= regs.reg[rs2_unsigned] {
			regs.pc = regs.pc + offset
		}
	default:
		return fmt.Errorf("invalid func3(val=%v) on BInstr", Instr.func3)
	}

	return nil
}

type UInstr struct {
	imm    uint32 // 20 bit offset
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

func (Inst UInstr) Execute(mem Memory, regs *Registers) error {
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
		return unknowOpcodeError(Inst.opcode, UInstrType)
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

func (Instr JInstr) Imm() uint32 {
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

	return signed
}

func (Instr JInstr) Execute(mem Memory, regs *Registers) error {
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
	} else {
		return unknowOpcodeError(Instr.opcode, JInstrType)
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

func CreateADDI(src int, dst int, imm uint32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_ADDI, opcode: OP_IMM}
}

func CreateSLLI(src int, dst int, imm uint32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SLLI, opcode: OP_IMM}
}

func CreateSLRI(src int, dst int, imm uint32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: FUNC3_SRLI, opcode: OP_IMM}
}

func CreateSRAI(src int, dst int, imm uint32) IInstr {
	if imm > 31 {
		panic("Invalid SRAI, the immediate should not be bigger then 31")
	}
	imm = imm + (32 << 5)
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
	return UInstr{imm: ReinterpreteAsUnsigned(imm), rd: dst, opcode: LUI}
}

func CreateAUIPC(imm int32, dst int32) UInstr {
	return UInstr{imm: ReinterpreteAsUnsigned(imm), rd: dst, opcode: AUIPC}
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
		log.Panic("imm in jal must be an even number")
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

func CreateJALR(imm uint32, link_reg int, addr_reg int) IInstr {
	return IInstr{imm: imm, rd: link_reg, opcode: JALR, rs1: addr_reg}
}

// Private function, for internal use only.
func createBRANCH(imm uint32, rs1 uint32, rs2 uint32, func3 uint32) BInstr {
	// split up the immediate
	// imm0: 1 bit -> 11
	// imm1: 4 bit -> 4:1
	// imm2: 6 bit -> 10:5
	// imm3: 1 bit -> 12
	imm0 := bitSliceBetween(imm, 11, 11)
	imm1 := bitSliceBetween(imm, 1, 4)
	imm2 := bitSliceBetween(imm, 5, 10)
	imm3 := bitSliceBetween(imm, 12, 12)

	return BInstr{imm3: imm3, imm2: imm2, rs2: rs2, rs1: rs1, func3: func3, imm1: imm1, imm0: imm0, opcode: BRANCH}
}

func CreateBEQ(imm uint32, rs1 uint32, rs2 uint32) BInstr {
	return createBRANCH(imm, rs1, rs2, FUNC3_BEQ)
}

func CreateBNE(imm uint32, rs1 uint32, rs2 uint32) BInstr {
	return createBRANCH(imm, rs1, rs2, FUNC3_BNE)
}

func CreateBGE(imm uint32, rs1 uint32, rs2 uint32) BInstr {
	return createBRANCH(imm, rs1, rs2, FUNC3_BGE)
}

func CreateBGEU(imm uint32, rs1 uint32, rs2 uint32) BInstr {
	return createBRANCH(imm, rs1, rs2, FUNC3_BGEU)
}

func CreateBLT(imm uint32, rs1 uint32, rs2 uint32) BInstr {
	return createBRANCH(imm, rs1, rs2, FUNC3_BLT)
}

func CreateBLTU(imm uint32, rs1 uint32, rs2 uint32) BInstr {
	return createBRANCH(imm, rs1, rs2, FUNC3_BLTU)
}

func CreateLoad(offset int32, addr int, func3 int8, dst int) IInstr {
	return IInstr{
		imm:    ReinterpreteAsUnsigned(offset),
		rs1:    addr,
		func3:  func3,
		rd:     dst,
		opcode: LOAD,
	}
}

func CreateLW(offset int32, addr int, dst int) IInstr {
	return CreateLoad(offset, addr, FUNC3_LW, dst)
}

func CreateLH(offset int32, addr int, dst int) IInstr {
	return CreateLoad(offset, addr, FUNC3_LH, dst)
}

func CreateLB(offset int32, addr int, dst int) IInstr {
	return CreateLoad(offset, addr, FUNC3_LB, dst)
}

func CreateStore(offset int32, src int, base int, func3 int8) SInstr {
	// imm1   uint32
	// rs2    uint32
	// rs1    uint32
	// func3  int8
	// imm0   uint32
	// opcode int8

	offset_unsigned := ReinterpreteAsUnsigned(offset)
	// imm0 -> offset[0:4]
	// imm1 -> offset[5:11]
	imm0 := (offset_unsigned << (32 - 5)) >> (32 - 5)
	imm1 := offset_unsigned >> 5

	return SInstr{
		imm1:   imm1,
		rs2:    src,
		rs1:    base,
		func3:  func3,
		imm0:   imm0,
		opcode: STORE,
	}
}

func CreateSB(offset int32, src int, base int) SInstr {
	return CreateStore(offset, src, base, FUNC3_SB)
}

func CreateSH(offset int32, src int, base int) SInstr {
	return CreateStore(offset, src, base, FUNC3_SH)
}

func CreateSW(offset int32, src int, base int) SInstr {
	return CreateStore(offset, src, base, FUNC3_SW)
}
