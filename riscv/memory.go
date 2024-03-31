package riscv

type Memory struct {
	data []int
}

func NewMemory(size int) Memory {
	return Memory{make([]int, size)}
}
