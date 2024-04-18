package main

import (
	"debug/elf"
	"emu/riscv"
	"flag"
	"fmt"
)

func PrintExecutableCodeSection(s *elf.Section, decoder *riscv.Decoder) {
	data, err := s.Data()
	if err != nil {
		panic("Invalid section passed to print function.")
	}

	// No compressed instructions right now, so we can assume all instructions are 4 bytes wide
	var cache [4]byte
	for i, b := range data {
		cache_i := i % 4
		cache[cache_i] = b
		if i > 0 && cache_i == 3 {
			// if we just read the last byte, fomat the instruction
			fmt.Printf("prog[%d]=%x%x%x%x \n", i-3, cache[3], cache[2], cache[1], cache[0])
			if decoder != nil {
				word := riscv.ByteArrayToWord(cache)
				instr, err := decoder.Decode(word)
				if err != nil {
					fmt.Printf("can't decode instruction with error: %v \n", err.Error())
					return
				}
				fmt.Printf("decoded as %v \n", instr)
			}
		}
	}
}

func main() {
	file := flag.String("file", "", "Elf file with risc machine code in it.")
	decodeInstr := flag.Bool("decode", true, "Decodes the instructions/")
	flag.Parse()

	if *file == "" {
		println("Please provide the --file argument.")
		return
	}

	f, err := elf.Open(*file)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Elf file with machine string=%s \n", f.Machine.String())
	section_index := int(-1)
	for i, section := range f.Sections {
		fmt.Printf("section[%d] has name-%s with sectionType=%s(%d) \n", i, section.Name, section.Type.String(), section.Type)
		if section.Type == elf.SHT_PROGBITS {
			if section_index > 0 {
				println("duplicate section type found")
			}
			section_index = i
		}
	}

	if section_index < 0 {
		panic("Error: executable section not found. \n")
	}
	fmt.Printf("Printing out section[%d] \n", section_index)

	var decoder *riscv.Decoder = nil
	if *decodeInstr {
		decoder = riscv.NewDecoder()
		// at the moment no extensions are supported
		decoder.RegisterBaseInstructionSet()
		fmt.Printf("Decoding instructions of base instruction set\n")
	}
	PrintExecutableCodeSection(f.Sections[section_index], decoder)

	f.Close()
}
