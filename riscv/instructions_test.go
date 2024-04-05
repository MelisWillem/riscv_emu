package riscv

import (
	"testing"
)

func CheckReg(regIndex int32, expected int32, r *Registers, t *testing.T) {
	if r.reg[regIndex] != expected {
		t.Logf("reg[%d]==%d and should be %d", regIndex, r.reg[1], expected)
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
