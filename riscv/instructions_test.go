package riscv

import (
	"testing"
)

func CheckReg(regIndex int, expected int32, r *Registers, t *testing.T) {
	if r.reg[regIndex] != expected {
		t.Logf("reg[%d]==%d and should be %d", regIndex, r.reg[regIndex], expected)
		t.Fail()
	}
}

func CheckPc(expected int32, r *Registers, t *testing.T) {
	if r.pc != expected {
		t.Logf("pc==%d and should be %d", r.pc, expected)
		t.Fail()
	}
}

func TestAddI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, add 2
	addi := CreateADDI(0, 1, 2)

	r.reg[0] = 4
	addi.Execute(&mem, &r)
	expected := int32(6)
	CheckReg(1, expected, &r, t)
}

func TestSLLI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, left shift of 2
	addi := CreateSLLI(0, 1, 2)

	r.reg[0] = 8
	expected := int32(32)
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

func TestSRLI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSLRI(0, 1, 2)

	r.reg[0] = 8
	expected := int32(2)
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

func TestSRLINegative(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSLRI(0, 1, 2)

	r.reg[0] = -32                // 11111111111111111111111111100000
	expected := int32(1073741816) // 00111111111111111111111111111000
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

func TestSRAI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSRAI(0, 1, 2)

	r.reg[0] = -32        // sext(11100000)
	expected := int32(-8) // sect(11111000)
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

// U-instr
func TestLui(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)
	// set 4096 (1 shifted by 12) in register 1
	I := CreateLui(1, 1)
	I.Execute(&mem, &r)

	expected := int32(4096)

	CheckReg(1, expected, &r, t)
}

func TestAUIPC(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)
	r.pc = 4
	// add 4096 (1 shifted by 12) to pc and put in register 1
	I := CreateAUIPC(1, 1)
	I.Execute(&mem, &r)

	expected := int32(4100)

	CheckReg(1, expected, &r, t)
}

func TestADD(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] + reg[3]
	r.reg[2] = 2
	r.reg[3] = 3
	expected := int32(5)

	I := CreateADD(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSUB(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] - reg[3]
	r.reg[2] = 2
	r.reg[3] = 3
	expected := int32(-1)

	I := CreateSUB(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSLT(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 4
	r.reg[3] = 5
	expected := int32(1) // 4<5==true

	I := CreateSLT(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSLTU(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 4
	r.reg[3] = -5
	expected := int32(1) // abs(4)<abs(-5)==true

	I := CreateSLTU(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestAND(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5         // 0101
	r.reg[3] = 4         // 0110
	expected := int32(4) // 0100

	I := CreateAND(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestOR(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5         // 0101
	r.reg[3] = 6         // 0110
	expected := int32(7) // 0111

	I := CreateOR(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestXOR(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5         // 0101
	r.reg[3] = 4         // 0110
	expected := int32(1) // 0011

	I := CreateXOR(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSLL(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5 // 0101
	r.reg[3] = 2
	expected := int32(20) // 010100

	I := CreateSLL(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSRA(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 20 // 010100
	r.reg[3] = 2
	expected := int32(5) // 0101

	I := CreateSRA(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSRL(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 20 // 010100
	r.reg[3] = 2
	expected := int32(5) // 0101

	I := CreateSRL(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestJAL(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	begin_pc := int32(10)
	r.pc = begin_pc
	pc_offset := int32(16)

	r.pc = 10 // we are at 10 -> should be in link register

	I := CreateJAL(pc_offset, reg_a0)
	// Assert(t, I.Imm(), pc_offset)
	I.Execute(&mem, &r)

	// make sure the link is saved
	// CheckReg(reg_a0, begin_pc+1, &r, t)
	CheckPc(begin_pc+pc_offset, &r, t)
}

func TestJALR(t *testing.T) {
	r := Registers{}
	mem := NewMemory(0)

	offset := int32(10)
	link_reg := reg_a0
	addr_reg := reg_a1
	begin_pc := int32(5)

	r.reg[addr_reg] = 15
	r.pc = begin_pc

	I := CreateJALR(offset, link_reg, addr_reg)
	I.Execute(&mem, &r)

	// make sure the link is saved
	CheckReg(link_reg, begin_pc+1, &r, t)
	// 10 + 15 = 25 -> setting least-sign to 0 results in 24
	CheckPc(24, &r, t)
}
