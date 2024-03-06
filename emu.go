package main

// import "fmt"
import "emu/riscv"

func main() {
	e := riscv.NewEmulator()
	e.Start()
	// fmt.Println("Hello world.")
}
