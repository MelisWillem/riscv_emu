package main

// import "fmt"
import "emu/riscv"

func main() {
	e := riscv.NewEmulator()
	program := make([]riscv.Instruction, 0)
	e.Start(program)
	// fmt.Println("Hello world.")
}
