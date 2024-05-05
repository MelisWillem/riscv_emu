package riscv

import (
	"testing"
)

func CheckReg(regIndex int, expected uint32, r Registers, t *testing.T) {
	if r.Reg(regIndex) != expected {
		t.Logf("reg[%d]==%d and should be %d", regIndex, r.Reg(regIndex), expected)
		t.Fail()
	}
}

func CheckPc(expected uint32, r Registers, t *testing.T) {
	if r.Pc() != expected {
		t.Logf("pc==%d and should be %d", r.Pc(), expected)
		t.Fail()
	}
}

func TestAddI(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(10)

	// read from x0, write to x1, add 2
	addi := CreateADDI(0, 1, 2)

	r.reg[0] = 4
	addi.Execute(&mem, &r)
	expected := uint32(6)
	CheckReg(1, expected, &r, t)
}

func TestSLLI(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(10)

	// read from x0, write to x1, left shift of 2
	addi := CreateSLLI(0, 1, 2)

	r.reg[0] = 8
	expected := uint32(32)
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

func TestSRLI(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSLRI(0, 1, 2)

	r.reg[0] = 8
	expected := uint32(2)
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

func TestSRLINegative(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSLRI(0, 1, 2)

	r.reg[0] = ReinterpreteAsUnsigned(-32) // 11111111111111111111111111100000
	expected := uint32(1073741816)         // 00111111111111111111111111111000
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

func TestSRAI(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSRAI(0, 1, 2)

	r.reg[0] = ReinterpreteAsUnsigned(-32) // sext(11100000)
	expected := ReinterpreteAsUnsigned(-8) // sect(11111000)
	addi.Execute(&mem, &r)
	CheckReg(1, expected, &r, t)
}

// U-instr
func TestLui(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)
	// set 4096 (1 shifted by 12) in register 1
	I := CreateLui(1, 1)
	I.Execute(&mem, &r)

	expected := uint32(4096)

	CheckReg(1, expected, &r, t)
}

func TestAUIPC(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)
	r.pc = 4
	// add 4096 (1 shifted by 12) to pc and put in register 1
	I := CreateAUIPC(1, 1)
	I.Execute(&mem, &r)

	expected := uint32(4100)

	CheckReg(1, expected, &r, t)
}

func TestADD(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] + reg[3]
	r.reg[2] = 2
	r.reg[3] = 3
	expected := uint32(5)

	I := CreateADD(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSUB(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] - reg[3]
	r.reg[2] = 2
	r.reg[3] = 3
	expected := ReinterpreteAsUnsigned(-1)

	I := CreateSUB(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSLT(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 4
	r.reg[3] = 5
	expected := uint32(1) // 4<5==true

	I := CreateSLT(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSLTU(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 4
	r.reg[3] = ReinterpreteAsUnsigned(-5)
	expected := uint32(1) // abs(4)<abs(-5)==true

	I := CreateSLTU(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestAND(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5          // 0101
	r.reg[3] = 4          // 0110
	expected := uint32(4) // 0100

	I := CreateAND(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestOR(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5          // 0101
	r.reg[3] = 6          // 0110
	expected := uint32(7) // 0111

	I := CreateOR(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestXOR(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5          // 0101
	r.reg[3] = 4          // 0110
	expected := uint32(1) // 0011

	I := CreateXOR(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSLL(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 5 // 0101
	r.reg[3] = 2
	expected := uint32(20) // 010100

	I := CreateSLL(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSRA(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 20 // 010100
	r.reg[3] = 2
	expected := uint32(5) // 0101

	I := CreateSRA(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestSRL(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	// reg[1] = reg[2] < reg[3]
	r.reg[2] = 20 // 010100
	r.reg[3] = 2
	expected := uint32(5) // 0101

	I := CreateSRL(1, 2, 3)
	I.Execute(&mem, &r)

	CheckReg(1, expected, &r, t)
}

func TestJAL(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	begin_pc := uint32(10)
	r.pc = begin_pc
	pc_offset := int32(16)

	r.pc = 10 // we are at 10 -> should be in link register

	I := CreateJAL(pc_offset, reg_a0)
	// Assert(t, I.Imm(), pc_offset)
	I.Execute(&mem, &r)

	// make sure the link is saved
	// CheckReg(reg_a0, begin_pc+1, &r, t)
	CheckPc(begin_pc+uint32(pc_offset), &r, t)
}

func TestJALR(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	offset := uint32(10)
	link_reg := reg_a0
	addr_reg := reg_a1
	begin_pc := uint32(5)

	r.reg[addr_reg] = 15
	r.pc = begin_pc

	I := CreateJALR(offset, link_reg, addr_reg)
	I.Execute(&mem, &r)

	// make sure the link registers points to the next instruction
	CheckReg(link_reg, begin_pc+4, &r, t)
	// 10 + 15 = 25 -> setting least-sign to 0 results in 24
	CheckPc(24, &r, t)
}

func TestBEQ(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(0)

	offset := ReinterpreteAsUnsigned(int32(10))
	begin_pc := uint32(5)

	r.reg[1] = 1
	r.reg[2] = 2
	r.reg[3] = 1

	r.pc = begin_pc

	I := CreateBEQ(offset, 1, 2)
	Assert(t, offset, I.imm())

	I.Execute(&mem, &r)

	// x1 != x2 -> should not branch
	CheckPc(begin_pc, &r, t)

	// x1 == x3 -> should branch
	I = CreateBEQ(offset, 1, 3)
	I.Execute(&mem, &r)

	CheckPc(begin_pc+10, &r, t)
}

func TestCreateStoreOffset(t *testing.T) {
	inputs := []int32{0, 31, 63}

	for _, offset := range inputs {
		I := CreateStore(offset, 0, 0, 0)

		if offset < 32 && I.imm() != I.imm0 {
			t.Logf("I.imm()=%d and I.imm0=%d", I.imm(), I.imm0)
			t.Fail()
		}

		if I.imm() != uint32(offset) {
			t.Logf("I.imm()!=offset I.Imm()=%d and offset=%d imm0=%d imm1=%d", I.imm(), offset, I.imm0, I.imm1)
			t.Fail()
		}
	}
}

func TestSWLW(t *testing.T) {
	r := RegistersImpl{}
	mem := NewMemory(20)

	offset := int32(10)
	addr := uint32(1)
	addrReg := 0
	LdReg := 1
	StReg := 2

	r.reg[addrReg] = addr
	wordToBeStored := uint32(5)
	r.SetReg(StReg, wordToBeStored) // the word to be store==5

	IStore := CreateSW(offset, StReg, addrReg)
	errStore := IStore.Execute(&mem, &r)
	if errStore != nil {
		t.Errorf("store instruction failed with error=%v", errStore.Error())
	}

	val, _ := mem.Load(1+uint32(offset), 4) // assume offset is not negative
	if val != wordToBeStored {
		t.Errorf("word not saved in memory, mem load results in %d but should be %d", val, wordToBeStored)
	}

	ILoad := CreateLW(offset, addrReg, LdReg)
	errLoad := ILoad.Execute(&mem, &r)
	if errLoad != nil {
		t.Errorf("load instruction failed with error=%v", errLoad.Error())
	}

	if r.Reg(LdReg) != wordToBeStored {
		t.Errorf("Load result(value=%d) != %d", r.Reg(LdReg), wordToBeStored)
	}
}
