package core

//CPURegisters - 6502  registers
type CPURegisters struct {
	Accumulator    uint8
	X              uint8
	Y              uint8
	StackPointer   uint8
	Status         uint8
	ProgramCounter uint16
}

//NewCPURegisters initializes all the registers
func NewCPURegisters() *CPURegisters {
	r := CPURegisters{
		Accumulator:    0x00,
		X:              0x00,
		Y:              0x00,
		StackPointer:   0xfd,
		Status:         UnusedBit,
		ProgramCounter: 0x0000,
	}
	return &r
}

//CPU represents the state of a 65c02.
type CPU struct {
	Registers      *CPURegisters
	Bus            SystemBus //system Bus
	operand        uint8     //operand for the current instruction
	operandAddress uint16    //address of the operand for the current instruction
	waiting        bool      //WAI instruction flag
	stopped        bool      //STP instruction flag
	handlingNMI    bool      //indicates whether the cpu is currently handling an NMI
	nmiQueue       int       //NMIs that occurred while handling other NMIs will increment this counter
}

//NewCPU returns an initialized CPU
func NewCPU(bus SystemBus, registers *CPURegisters) *CPU {
	c := CPU{
		Registers:      registers,
		Bus:            bus,
		waiting:        false,
		stopped:        false,
		handlingNMI:    false,
		operand:        0x00,
		operandAddress: 0x0000,
		nmiQueue:       0,
	}
	return &c
}

//Execute one instruction
func (cpu *CPU) Execute() {
	if !cpu.stopped && !cpu.waiting {
		opcode := cpu.read(cpu.Registers.ProgramCounter)
		cpu.Registers.ProgramCounter++
		instr := instructionTable[opcode]
		if instr.addressing != nil {
			instr.addressing(cpu)
		}
		if instr.operation != nil {
			instr.operation(cpu)
		}
		cpu.operand = 0x00
		cpu.operandAddress = 0x0000
	}
}

func (cpu *CPU) interrupt(vectorLowByte uint16) {
	pch := uint8((cpu.Registers.ProgramCounter & 0xff) >> 8)
	pcl := uint8(cpu.Registers.ProgramCounter & 0xff)
	cpu.pushStack(pch)
	cpu.pushStack(pcl)
	cpu.pushStack(cpu.Registers.Status)
	cpu.Registers.ProgramCounter = vectorLowByte
	cpu.abs()
	cpu.Registers.ProgramCounter = cpu.operandAddress
}

//Interrupt sends a maskable hardware interrupt
func (cpu *CPU) Interrupt() {
	if !cpu.stopped && !cpu.testStatusBit(InterruptDisableBit) {
		if cpu.waiting {
			cpu.waiting = false
		}
		cpu.setStatusBit(InterruptDisableBit, true)
		cpu.interrupt(vectorBRKL)
	}
}

//NMInterrupt sends a non-maskable hardware interrupt
func (cpu *CPU) NMInterrupt() {
	if cpu.handlingNMI {
		cpu.nmiQueue++
	} else {
		if cpu.waiting {
			cpu.waiting = false
		}
		if cpu.nmiQueue > 0 {
			cpu.nmiQueue--
		}
		cpu.interrupt(vectorNMIBL)
	}
}

//Reset - resets the cpu to a known state
func (cpu *CPU) Reset() {
	cpu.Registers = NewCPURegisters()
	cpu.waiting = false
	cpu.stopped = false
	cpu.handlingNMI = false
	cpu.nmiQueue = 0
	cpu.Registers.ProgramCounter = vectorRESBL
	cpu.abs()
	cpu.Registers.ProgramCounter = cpu.operandAddress
}
