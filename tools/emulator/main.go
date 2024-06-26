package main

// import "fmt"
import (
	"debug/elf"
	"emu/riscv"
	"flag"
	"log"
)

func processExecutionBody(mem riscv.Memory, r riscv.Registers, executableData []byte, decoder *riscv.Decoder) {
	var cache [4]byte
	for i, b := range executableData {
		cache_i := i % 4
		cache[cache_i] = b
		if i > 0 && cache_i == 3 {
			// if we just read the last byte, fomat the instruction
			if decoder != nil {
				word := riscv.ByteArrayToWord(cache)
				log.Printf("decoding instruction at byte offset %v", i-3)
				instr, err := decoder.Decode(word)
				if err != nil {
					log.Fatalf("can't decode instruction with error: %v", err.Error())
				}
				log.Printf("executing instruction I=%s", instr.String())
				err = instr.Execute(mem, r)
				if err != nil {
					log.Fatalf("can't execute instruction (%s) with error: %v", instr.String(), err.Error())
				}
			}
		}
	}
}

func main() {
	file := flag.String("file", "", "Elf file with risc machine code in it.")
	logRegisterChanged := flag.Bool("log-reg", false, "Log all register changes")
	memory_size := flag.Int("memory_size", 100, "Size of the memory")
	memory_offset := flag.Int("memory_offset", 0, "Begin address of memory")
	flag.Parse()

	if *file == "" {
		println("Please provide the --file argument.")
		return
	}

	f, err := elf.Open(*file)
	if err != nil {
		log.Panic(err.Error())
	}

	mem := riscv.NewMemoryWithOffset(*memory_size, uint32(*memory_offset))
	// var r *riscv.RegistersImpl = &riscv.RegistersImpl{}
	var r riscv.Registers = &riscv.RegistersImpl{}
	if *logRegisterChanged {
		r = riscv.NewLoggedRegisters(r)
	}

	decoder := riscv.NewDecoder()
	log.Println("Registering base instruction set in decoder")
	decoder.RegisterBaseInstructionSet()

	for i, section := range f.Sections {
		if section.Type == elf.SHT_PROGBITS {
			log.Printf("Executing body of section %d with name %s \n", i, section.Name)
			data, err := section.Data()
			if err != nil {
				log.Panicf("Invalid section passed to print function.")
			}

			processExecutionBody(&mem, r, data, decoder)
		}
	}

}
