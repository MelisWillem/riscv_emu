package riscv

import (
	"fmt"
)

type emulator struct {
	programy []int32
	data     []int32
	regs     Registers
	pcu      Pcu
	fetcher  Fetcher
	decoder  Decoder
}

func (e emulator) Start() {
	fmt.Println("Starting emulator...")

	raw_inst, cont := e.fetcher.Fetch(e.regs.pc)
	for cont {
		inst := e.decoder.Decode(raw_inst)
		inst.Execute(&e.pcu)
		e.pcu.NextStep(&e.regs)
		// todo memory writeback
	}
}

func NewEmulator() *emulator {
	a := new(emulator)
	return a
}
