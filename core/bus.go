package core

//Maximum addressable range for the wdc65c02
const MaxBusSize int = 1024 * 64

//Implement this interface if you want to use a custom Bus.
type SystemBus interface {
	//SystemBus.Read returns the value read from device at addr.
	Read(addr uint16) uint8
	//SystemBus.Write writes value to device located at addr. Error
	//should be returned if attempting to write to a readonly device.
	Write(addr uint16, val uint8) error
}

//BasicBus represents the simplest Bus that could be used with the wdc65c02 cpu.
//The only device connected is 64K of RAM.
type BasicBus struct {
	memory []uint8
}

func NewBasicBus() *BasicBus {
	return &BasicBus{
		make([]uint8, MaxBusSize),
	}
}

func (bus *BasicBus) Read(addr uint16) uint8 {
	return bus.memory[addr]
}

func (bus *BasicBus) Write(addr uint16, val uint8) error {
	bus.memory[addr] = val
	return nil
}
