package riscv

import (
	"testing"
)

func TestStoreLoadByte(t *testing.T) {
	dataToStore := uint32(256 + 1) // 100000001
	// the bit of 256 should not be stored as we will
	// store just one byte
	expectedLoadValue := uint32(1)
	mem := NewMemory(10)

	for addr := uint32(0); addr < 4; addr++ {

		mem.StoreByte(addr, dataToStore)
		loadRes, err := mem.LoadByte(addr)

		if err != nil {
			t.Logf("Failed to load at addr=%d with erro %v", addr, err)
			t.Fail()
		}

		if loadRes != expectedLoadValue {
			t.Logf("Invalid load of addr=%v value(val=%v) != expected(val=%v) mem.data[0]=%v", addr, loadRes, expectedLoadValue, mem.data[0])
			t.Fail()
		}

		if expectedLoadValue != uint32(mem.data[addr]) {
			t.Logf("Invalid value in memory expected=%d but value=%d", expectedLoadValue, mem.data[addr])
			t.Fail()
		}
	}
}

func TestStoreLoad(t *testing.T) {
	dataToStore := uint32(256 + 1) // 100000001
	expectedLoadValue := dataToStore
	mem := NewMemory(10)

	for addr := uint32(0); addr < 4; addr++ {
		errStore := mem.Store(addr, dataToStore, 2)
		if errStore != nil {
			t.Logf("Failed to store at addr=%d with error %v", addr, errStore)
			t.Fail()
		}

		loadRes, errLoad := mem.Load(addr, 2)
		if errLoad != nil {
			t.Logf("Failed to load at addr=%d with error %v", addr, errLoad)
			t.Fail()
		}

		if loadRes != expectedLoadValue {
			t.Logf("Invalid load of addr=%v value(val=%v) != expected(val=%v) mem.data[0]=%v", addr, loadRes, expectedLoadValue, mem.data[0])
			t.Fail()
		}
	}
}
