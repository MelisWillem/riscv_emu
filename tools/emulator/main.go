package main

// import "fmt"
import (
	"debug/elf"
	"emu/riscv"
	"flag"
	"fmt"
)

func processExecutionBody(mem *riscv.Memory, r *riscv.Registers, executableData []byte) {
	for i := 0; i < len(executableData); i = i + 4 {
		registration_map := map[int8]int8{} // parse the 4 bytes left to right
		riscv.RegisterIInstr(&registration_map)
		riscv.RegisterRInstr(&registration_map)
		word := uint32(0)
		for j := 0; j < 4; j++ {
			word = word | (uint32(executableData[i+j]) << (j * 8))

			opcode := executableData[i] & uint8(63)
			instrType, isRegistered := registration_map[int8(opcode)]
			if !isRegistered {
				panic(fmt.Sprintf("Unregistered opcode: %d", opcode))
			}

			var I riscv.Instruction
			switch instrType {
			case riscv.IInstrType:
				I = riscv.DecodeIInstr(word)
			case riscv.RInstrType:
				I = riscv.DecodeRInstr(word)
			case riscv.UInstrType:
				I = riscv.DecodeUInstr(word)
			case riscv.JInstrType:
				I = riscv.DecodeJInstr(word)
			default:
				fmt.Printf("Cannot decode instr %d", word)
			}
			I.Execute(mem, r)

		}
		fmt.Printf("word=%d \n", word)
		// execute the instruction
	}
}

func main() {
	file := flag.String("file", "", "Elf file with risc machine code in it.")
	flag.Parse()

	if *file == "" {
		println("Please provide the --file argument.")
		return
	}

	f, err := elf.Open(*file)
	if err != nil {
		panic(err.Error())
	}

	mem := riscv.NewMemory(100)
	r := riscv.Registers{}

	for i, section := range f.Sections {
		if section.Type == elf.SHT_PROGBITS {
			fmt.Printf("Executing bod of section %d with name %s \n", i, section.Name)
			data, err := section.Data()
			if err != nil {
				panic("Invalid section passed to print function.")
			}

			processExecutionBody(&mem, &r, data)
		}
	}

}
