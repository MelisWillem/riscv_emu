package riscv

import (
	"fmt"
	"log"
)

type Memory interface {
	StoreByte(addr uint32, data uint32) error
	Store(addr uint32, data uint32, numBytes uint32) error
	LoadByte(addr uint32) (uint32, error)
	Load(addr uint32, numBytes uint32) (uint32, error)
	Len() int
}

type MemoryImpl struct {
	// the memory is byte addressed
	data []uint8
}

func (mem *MemoryImpl) StoreByte(addr uint32, data uint32) error {
	if addr >= uint32(mem.Len()) {
		return fmt.Errorf("out of range error on store max addr=%v but actual addr=%v", mem.Len(), addr)
	}

	filteredData := data & uint32(255)
	dataByte := uint8(data & filteredData)
	mem.data[addr] = dataByte

	return nil
}

func (mem *MemoryImpl) Store(addr uint32, data uint32, numBytes uint32) error {
	// numBitsToSkip := (4 - numBytes) * 8
	// data = (data << numBitsToSkip) >> numBitsToSkip
	if numBytes > 4 && numBytes < 1 {
		return fmt.Errorf("numBytes must be (0 < numBytes <= 4) but is %d", numBytes)
	}

	for i := uint32(0); i < numBytes; i++ {
		mem.StoreByte(addr+i, data)
		data = data >> 1
	}

	return nil
}

func (mem *MemoryImpl) LoadByte(addr uint32) (uint32, error) {
	if addr >= uint32(mem.Len()) {
		return 0, fmt.Errorf("out of range error on load max addr=%v but actual addr=%v", mem.Len(), addr)
	}

	return uint32(mem.data[addr]), nil
}

func (mem *MemoryImpl) Load(addr uint32, numBytes uint32) (uint32, error) {
	if numBytes > 4 && numBytes < 1 {
		return 0, fmt.Errorf("numBytes must be (0 < numBytes <= 4) but is %d", numBytes)
	}
	data := uint32(0)
	for i := uint32(0); i < numBytes; i++ {
		byteData, err := mem.LoadByte(addr + i)
		if err != nil {
			return 0, err
		}
		data |= (byteData << i)
	}

	return data, nil
}

func (mem *MemoryImpl) Len() int {
	return len(mem.data)
}

func NewMemory(size int) MemoryImpl {
	return MemoryImpl{make([]uint8, size)}
}

type LoggedMemory struct {
	mem MemoryImpl
}

func (mem *LoggedMemory) StoreByte(addr uint32, data uint32) error {
	log.Printf("StoreByte addr=%d data=%d", addr, data)
	return mem.mem.StoreByte(addr, data)
}

func (mem *LoggedMemory) Store(addr uint32, data uint32, numBytes uint32) error {
	log.Printf("Store addr=%d data=%d numBytes=%d", addr, data, numBytes)
	return mem.mem.Store(addr, data, numBytes)
}

func (mem *LoggedMemory) LoadByte(addr uint32) (uint32, error) {
	log.Printf("LoadByte addr=%d", addr)
	return mem.mem.LoadByte(addr)
}

func (mem *LoggedMemory) Load(addr uint32, numBytes uint32) (uint32, error) {
	log.Printf("Load addr=%d numBtes=%d", addr, numBytes)
	return mem.mem.Load(addr, numBytes)
}

func (mem *LoggedMemory) Len() int {
	return mem.mem.Len()
}
