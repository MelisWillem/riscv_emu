package main

// import "fmt"
import (
	"debug/elf"
	"emu/riscv"
	"flag"
	"fmt"
)

func processExecutionBody(mem *riscv.Memory, r *riscv.Registers, executableData []byte, decoder *riscv.Decoder) error {
	var cache [4]byte
	for i, b := range executableData {
		cache_i := i % 4
		cache[cache_i] = b
		if i > 0 && cache_i == 3 {
			// if we just read the last byte, fomat the instruction
			if decoder != nil {
				word := riscv.ByteArrayToWord(cache)
				instr, err := decoder.Decode(word)
				if err != nil {
					return fmt.Errorf("can't decode instruction with error: %v", err.Error())
				}
				instr.Execute(mem, r)
			}
		}
	}
	return nil
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
	decoder := riscv.NewDecoder()
	decoder.RegisterBaseInstructionSet()

	for i, section := range f.Sections {
		if section.Type == elf.SHT_PROGBITS {
			fmt.Printf("Executing bod of section %d with name %s \n", i, section.Name)
			data, err := section.Data()
			if err != nil {
				panic("Invalid section passed to print function.")
			}

			err = processExecutionBody(&mem, &r, data, decoder)
			if err != nil {
				fmt.Printf("Failed to process executable data with error: %v", err.Error())
			}
		}
	}

}
