package riscv

type Decoder struct {
}

func (d Decoder) Decode(input int32) Instruction {
	return InvalidInstrction{}
}
