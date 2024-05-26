package riscv

import "log"

type RegistersImpl struct {
	reg [32]uint32
	pc  uint32
}

type Registers interface {
	Reg(i int) uint32
	SetReg(i int, data uint32)

	Pc() uint32
	SetPc(uint32)
}

func (r *RegistersImpl) Reg(i int) uint32 {
	return r.reg[i]
}

func (r *RegistersImpl) SetReg(i int, data uint32) {
	r.reg[i] = data
}

func (r *RegistersImpl) Pc() uint32 {
	return r.pc
}

func (r *RegistersImpl) SetPc(pc uint32) {
	r.pc = pc
}

const (
	reg_zero int = 0
	reg_ra   int = 1
	reg_sp   int = 2
	reg_gp   int = 3
	reg_tp   int = 4
	reg_t0   int = 5
	reg_t1   int = 6
	reg_t2   int = 7
	reg_s0   int = 8
	reg_fp   int = 8
	reg_s1   int = 9
	reg_a0   int = 10
	reg_a1   int = 11
	reg_a2   int = 12
	reg_a3   int = 13
	reg_a4   int = 14
	reg_a5   int = 15
	reg_a6   int = 16
	reg_a7   int = 17
	reg_s2   int = 18
	reg_s3   int = 19
	reg_s4   int = 20
	reg_s5   int = 21
	reg_s6   int = 22
	reg_s7   int = 23
	reg_s8   int = 24
	reg_s9   int = 25
	reg_s10  int = 26
	reg_s11  int = 27
	reg_t3   int = 28
	reg_t4   int = 29
	reg_t5   int = 30
	reg_t6   int = 31
)

type LoggedRegisters struct {
	reg Registers
}

func (r *LoggedRegisters) Reg(i int) uint32 {
	return r.reg.Reg(i)
}

func (r *LoggedRegisters) SetReg(i int, data uint32) {
	log.Printf("Setting reg[%d]=%d", i, data)
	r.reg.SetReg(i, data)
}

func (r *LoggedRegisters) Pc() uint32 {
	return r.reg.Pc()
}

func (r *LoggedRegisters) SetPc(pc uint32) {
	log.Printf("Setting pc=%d", pc)
	r.reg.SetPc(pc)
}

func NewLoggedRegisters(r Registers) *LoggedRegisters {
	return &LoggedRegisters{r}
}
