package riscv

import (
	"testing"
)

func TestAddI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, add 2
	addi := CreateADDI(0, 1, 2)

	r.reg[0] = 4
	addi.Execute(&mem, &r)
	if r.reg[1] != 6 {
		t.Fail()
	}
}

func TestSLLI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, left shift of 2
	addi := CreateSLLI(0, 1, 2)

	r.reg[0] = 8
	expected := int32(32)
	addi.Execute(&mem, &r)
	if r.reg[1] != expected {
		t.Logf("reg[1]==%d and should be %d", r.reg[1], expected)
		t.Fail()
	}
}

func TestSRLI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSLRI(0, 1, 2)

	r.reg[0] = 8
	expected := int32(2)
	addi.Execute(&mem, &r)
	if r.reg[1] != expected {
		t.Logf("reg[1]==%d and should be %d", r.reg[1], expected)
		t.Fail()
	}
}

func TestSRLINegative(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSLRI(0, 1, 2)

	r.reg[0] = -32                // 11111111111111111111111111100000
	expected := int32(1073741816) // 00111111111111111111111111111000
	addi.Execute(&mem, &r)
	if r.reg[1] != expected {
		t.Logf("reg[1]==%d and should be %d", r.reg[1], expected)
		t.Fail()
	}
}

func TestSRAI(t *testing.T) {
	r := Registers{}
	mem := NewMemory(10)

	// read from x0, write to x1, right shift of 2
	addi := CreateSRAI(0, 1, 2)

	r.reg[0] = -32        // sext(11100000)
	expected := int32(-8) // sect(11111000)
	addi.Execute(&mem, &r)
	if r.reg[1] != expected {
		t.Logf("reg[1]==%d and should be %d", r.reg[1], expected)
		t.Fail()
	}
}
