package riscv

import (
	"fmt"
)

type emulator struct {
	mem  Memory
	regs Registers
}

func (e emulator) Start(program []Instruction) {
	fmt.Println("Starting emulator...")
	for index, inst := range program {
		err := inst.Execute(&e.mem, &e.regs)
		if err != nil {
			fmt.Printf("Error at instruction %d with error=%s", index, err.Error())
		}
	}
}

func NewEmulator() *emulator {
	a := new(emulator)
	return a
}
