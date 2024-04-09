package main

import (
	"debug/elf"
	"flag"
	"fmt"
)

func PrintExecutableCodeSection(s *elf.Section) {
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
		}
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
	PrintExecutableCodeSection(f.Sections[section_index])

	f.Close()
}
