package riscv

import "fmt"

type Memory struct {
	// the memory is byte addressed
	data []uint8
}

func (mem Memory) StoreByte(addr uint32, data uint32) error {
	if addr >= uint32(mem.Len()) {
		return fmt.Errorf("out of range error on store max addr=%v but actual addr=%v", mem.Len(), addr)
	}

	filteredData := data & uint32(255)
	dataByte := uint8(data & filteredData)
	mem.data[addr] = dataByte

	return nil
}

func (mem Memory) Store(addr uint32, data uint32, numBytes uint32) error {
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

func (mem Memory) LoadByte(addr uint32) (uint32, error) {
	if addr >= uint32(mem.Len()) {
		return 0, fmt.Errorf("out of range error on load max addr=%v but actual addr=%v", mem.Len(), addr)
	}

	return uint32(mem.data[addr]), nil
}

func (mem Memory) Load(addr uint32, numBytes uint32) (uint32, error) {
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

func (mem Memory) Len() int {
	return len(mem.data)
}

func NewMemory(size int) Memory {
	return Memory{make([]uint8, size)}
}
