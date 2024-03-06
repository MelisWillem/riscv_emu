package riscv

const (
	R    int = 1
	I    int = 2
	S    int = 3
	U    int = 4
	IImm int = 5
)

const (
	ADDI int = 0
	SLTI int = 1
	ANDI int = 2
	ORI  int = 3
	XORI int = 4
	SLLI int = 5
	SRLI int = 6
	SRAI int = 7
)

type Instruction interface {
	Execute(pcu *Pcu)
	// print()
}

type InvalidInstrction struct {
}

func (Inst InvalidInstrction) Execute(pcu *Pcu) {
	panic("Trying to execute invalid instruction...")
}

type RInstr struct {
	rs2 int
	rs1 int
	rd  int
}

type IInstr struct {
	imm    int32
	rs1    int
	rd     int
	func3  int
	opcode int
}

func (Inst IInstr) Execute(pcu *Pcu, regs *Registers) {
	switch Inst.opcode {
	case OP_IMM:
		op_imm_execute(Inst, pcu, regs)
		break
	default:
		panic("Unknown operator type on IInstr")
	}
}

func op_imm_execute(Inst IInstr, pcu *Pcu, regs *Registers) {
	switch Inst.func3 {
	case ADDI:
		regs.reg[Inst.rs1] = regs.reg[(Inst.rd)] + Inst.imm
		break
	default:
		panic("Invalid op_imm instruction.")

	}
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
	imm1 int
	rd   int
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
)

func CreateADDI(src int, dst int, imm int32) IInstr {
	return IInstr{rs1: dst, rd: src, imm: imm, func3: ADDI, opcode: OP_IMM}
}

func CreateMV(src int, dst int) IInstr {
	return CreateADDI(src, dst, 0)
}

func Nop() IInstr {
	// ADDI x0, x0, 0
	return CreateADDI(0, 0, 0)
}

func Ld(base int, width int, dest int, offset int32) IInstr {
	return IInstr{imm: offset, rs1: base, func3: width, rd: dest, opcode: LOAD}
}
