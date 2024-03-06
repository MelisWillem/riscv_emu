package riscv

import (
	"testing"
)

func TestAddI(t *testing.T) {
	r := Registers{}
	pcu := Pcu{}

	// read from x0, write to x1, add 2
	addi := CreateADDI(0, 1, 2)

	r.reg[0] = 4
	addi.Execute(&pcu, &r)
	if r.reg[1] != 6 {
		t.Fail()
	}
}
