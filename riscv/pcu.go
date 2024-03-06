package riscv

type Pcu struct {
}

func (pcu Pcu) NextStep(regs *Registers) {
	regs.pc = regs.pc + 1
}
