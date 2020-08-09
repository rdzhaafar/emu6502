package core

const (
	vectorBRKH  uint16 = 0xfffe
	vectorBRKL  uint16 = 0xffff
	vectorRESBH uint16 = 0xfffc
	vectorRESBL uint16 = 0xfffd
	vectorNMIBH uint16 = 0xfffa
	vectorNMIBL uint16 = 0xfffb

	bit7 uint8 = 0x80
	bit6 uint8 = 0x40
	bit5 uint8 = 0x20
	bit4 uint8 = 0x10
	bit3 uint8 = 0x08
	bit2 uint8 = 0x04
	bit1 uint8 = 0x02
	bit0 uint8 = 0x01

	//NegativeBit -> status register (p)
	NegativeBit = bit7
	//OverflowBit -> status register (p)
	OverflowBit = bit6
	//UnusedBit -> status register (p)
	UnusedBit = bit5
	//BreakBit -> status register (p)
	BreakBit = bit4
	//DecimalBit -> status register (p)
	DecimalBit = bit3
	//InterruptDisableBit -> status register (p)
	InterruptDisableBit = bit2
	//ZeroBit -> status register (p)
	ZeroBit = bit1
	//CarryBit -> status register (p)
	CarryBit = bit0
)

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
System Bus I/O wrappers
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func (cpu *CPU) read(addr uint16) uint8 {
	return cpu.Bus.Read(addr)
}

func (cpu *CPU) write(addr uint16, val uint8) {
	err := cpu.Bus.Write(addr, val)
	if err != nil {
		panic(err) //TODO: propagate the error to the appropriate handler
	}
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Status register helpers
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func (cpu *CPU) setStatusBit(flag uint8, val bool) {
	if val {
		cpu.Registers.Status |= flag
	} else {
		cpu.Registers.Status &= ^flag
	}
}

func (cpu *CPU) testStatusBit(flag uint8) bool {
	if cpu.Registers.Status&flag != 0 {
		return true
	}
	return false
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Addressing mode handlers
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

//absolute
func (cpu *CPU) abs() {
	cpu.operandAddress = uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.read(cpu.Registers.ProgramCounter)) << 8
	cpu.Registers.ProgramCounter++
	cpu.operand = cpu.read(cpu.operandAddress)
}

//absolute indexed indirect
func (cpu *CPU) aii() {
	indirectBaseAddress := uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	indirectBaseAddress += uint16(cpu.read(cpu.Registers.ProgramCounter)) << 8
	indirectBaseAddress += uint16(cpu.Registers.X)
	indirectAddress := uint16(cpu.read(indirectBaseAddress))
	indirectBaseAddress++
	indirectAddress += uint16(cpu.read(indirectBaseAddress)) << 8
	cpu.Registers.ProgramCounter = indirectAddress
	cpu.operandAddress = 0x0000
	cpu.operand = 0x00
}

//absolute indexed with X
func (cpu *CPU) aix() {
	cpu.operandAddress = uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.read(cpu.Registers.ProgramCounter)) << 8
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.Registers.X)
	cpu.operand = cpu.read(cpu.operandAddress)
}

//absolute indexed with Y
func (cpu *CPU) aiy() {
	cpu.operandAddress = uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.read(cpu.Registers.ProgramCounter)) << 8
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.Registers.Y)
	cpu.operand = cpu.read(cpu.operandAddress)
}

//absolute indirect
func (cpu *CPU) ai() {
	indirectAddress := uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	indirectAddress += uint16(cpu.read(cpu.Registers.ProgramCounter)) << 8
	cpu.Registers.ProgramCounter = indirectAddress
	cpu.operandAddress = 0x0000
	cpu.operand = 0x00
}

//immediate
func (cpu *CPU) imm() {
	cpu.operandAddress = cpu.Registers.ProgramCounter
	cpu.Registers.ProgramCounter++
	cpu.operand = cpu.read(cpu.operandAddress)
}

//program counter relative
func (cpu *CPU) pcr() {
	cpu.Registers.ProgramCounter += uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.operandAddress = 0x0000
	cpu.operand = 0x00
}

//zero page
func (cpu *CPU) zp() {
	cpu.operandAddress = uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operand = cpu.read(cpu.operandAddress)
}

//zero page indexed indirect
func (cpu *CPU) zpii() {
	baseAddress := uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	baseAddress += uint16(cpu.Registers.X)
	indirectAddress := uint16(cpu.read(baseAddress))
	cpu.operandAddress = uint16(cpu.read(indirectAddress))
	indirectAddress++
	cpu.operandAddress += uint16(cpu.read(indirectAddress)) << 8
	cpu.operand = cpu.read(cpu.operandAddress)
}

//zero page indexed with X
func (cpu *CPU) zpx() {
	cpu.operandAddress = uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.Registers.X)
	cpu.operand = cpu.read(cpu.operandAddress)
}

//zero page indexed with Y
func (cpu *CPU) zpy() {
	cpu.operandAddress = uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operandAddress += uint16(cpu.Registers.Y)
	cpu.operand = cpu.read(cpu.operandAddress)
}

//zero page indirect
func (cpu *CPU) zpi() {
	indirectAddress := uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	cpu.operandAddress = uint16(cpu.read(indirectAddress))
	indirectAddress++
	cpu.operandAddress += uint16(cpu.read(indirectAddress)) << 8
	cpu.operand = cpu.read(cpu.operandAddress)
}

//zero page indirect indexed with Y
func (cpu *CPU) zpiy() {
	indirectBaseAddress := uint16(cpu.read(cpu.Registers.ProgramCounter))
	cpu.Registers.ProgramCounter++
	indirectBaseAddress += uint16(cpu.Registers.Y)
	cpu.operandAddress = uint16(cpu.read(indirectBaseAddress))
	indirectBaseAddress++
	cpu.operandAddress += uint16(cpu.read(indirectBaseAddress)) << 8
	cpu.operand = cpu.read(cpu.operandAddress)
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Stack operations
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func (cpu *CPU) pullStack() uint8 {
	cpu.Registers.StackPointer++
	stackAddr := 0x0100 + uint16(cpu.Registers.StackPointer)
	return cpu.read(stackAddr)
}

func (cpu *CPU) pushStack(val uint8) {
	stackAddr := 0x0100 + uint16(cpu.Registers.StackPointer)
	cpu.write(stackAddr, val)
	cpu.Registers.StackPointer--
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Instruction handlers
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func (cpu *CPU) adc() {
	carry := uint16(0)
	if cpu.testStatusBit(CarryBit) {
		carry++
	}
	if cpu.testStatusBit(DecimalBit) {
		tmpl := uint16(cpu.Registers.Accumulator&0xf) + uint16(cpu.operand&0xf) + carry
		tmph := uint16(cpu.Registers.Accumulator&0xf0) + uint16(cpu.operand&0xf0)
		if tmpl > 0x9 {
			tmph += 0x10
			tmpl += 0x06
		}
		cpu.setStatusBit(
			OverflowBit,
			(^(uint16(cpu.Registers.Accumulator^cpu.operand))&(uint16(cpu.Registers.Accumulator)^tmph)&0x80) != 0,
		)
		if tmph > 0x90 {
			tmph += 0x60
		}
		cpu.setStatusBit(CarryBit, (tmph&0xff00) != 0)
		res := (tmpl & 0xf) | (tmph & 0xf0)
		cpu.setStatusBit(NegativeBit, res&0x80 != 0)
		cpu.setStatusBit(ZeroBit, res&0xff == 0)
		cpu.Registers.Accumulator = uint8(res)
	} else {
		res := uint16(cpu.Registers.Accumulator) + uint16(cpu.operand) + carry
		cpu.setStatusBit(
			OverflowBit,
			!(^((cpu.Registers.Accumulator&cpu.operand)&0x80) != 0) && (((uint16(cpu.Registers.Accumulator)^res)&0x80) != 0),
		)
		cpu.setStatusBit(ZeroBit, res&0xff == 0)
		cpu.setStatusBit(NegativeBit, res&0xff&0x80 != 0)
		cpu.Registers.Accumulator = uint8(res)
	}
}

func (cpu *CPU) sbc() {
	carry := uint16(0)
	if cpu.testStatusBit(CarryBit) {
		carry++
	}
	tmp := uint16(cpu.Registers.Accumulator) - uint16(cpu.operand) + carry - 1
	cpu.setStatusBit(
		OverflowBit,
		(((uint16(cpu.Registers.Accumulator)^tmp)&0x80) != 0) && (((cpu.Registers.Accumulator^cpu.operand)&0x80) != 0),
	)
	if cpu.testStatusBit(DecimalBit) {
		tmpl := uint16(cpu.Registers.Accumulator&0xf) - uint16(cpu.operand&0xf) + carry - 1
		if tmp > 0xff {
			tmp -= 0x60
		}
		if tmpl > 0xff {
			tmpl -= 0x6
		}
	}
	cpu.setStatusBit(CarryBit, (uint16(cpu.Registers.Accumulator)+carry-1) >= uint16(cpu.operand))
	cpu.setStatusBit(NegativeBit, tmp&0x80 != 0)
	cpu.setStatusBit(ZeroBit, tmp&0xff == 0)
	cpu.Registers.Accumulator = uint8(tmp)
}

func (cpu *CPU) and() {
	cpu.Registers.Accumulator &= cpu.operand
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) asl() {
	cpu.setStatusBit(CarryBit, cpu.operand&0x80 != 0)
	cpu.operand = cpu.operand << 1
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.setStatusBit(NegativeBit, cpu.operand&0x80 != 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

//ASL Accumulator
func (cpu *CPU) asla() {
	cpu.setStatusBit(CarryBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.Registers.Accumulator = cpu.Registers.Accumulator << 1
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
}

func (cpu *CPU) bcc() {
	if !cpu.testStatusBit(CarryBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bcs() {
	if cpu.testStatusBit(CarryBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) beq() {
	if cpu.testStatusBit(ZeroBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bmi() {
	if cpu.testStatusBit(NegativeBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bne() {
	if !cpu.testStatusBit(ZeroBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bpl() {
	if !cpu.testStatusBit(NegativeBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bvc() {
	if !cpu.testStatusBit(OverflowBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bvs() {
	if cpu.testStatusBit(OverflowBit) {
		cpu.pcr()
	}
}

func (cpu *CPU) bit() {
	cpu.setStatusBit(NegativeBit, cpu.operand&0x80 != 0)
	cpu.setStatusBit(OverflowBit, cpu.operand&0x40 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator&cpu.operand == 0)
}

func (cpu *CPU) bbr0() {
	if cpu.operand&^bit0 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr1() {
	if cpu.operand&^bit1 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr2() {
	if cpu.operand&^bit2 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr3() {
	if cpu.operand&^bit3 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr4() {
	if cpu.operand&^bit4 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr5() {
	if cpu.operand&^bit5 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr6() {
	if cpu.operand&^bit6 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbr7() {
	if cpu.operand&^bit7 == 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs0() {
	if cpu.operand&bit0 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs1() {
	if cpu.operand&bit1 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs2() {
	if cpu.operand&bit2 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs3() {
	if cpu.operand&bit3 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs4() {
	if cpu.operand&bit4 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs5() {
	if cpu.operand&bit5 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs6() {
	if cpu.operand&bit6 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) bbs7() {
	if cpu.operand&bit7 != 0 {
		cpu.pcr()
	}
}

func (cpu *CPU) rmb0() {
	cpu.operand &= ^bit0
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb1() {
	cpu.operand &= ^bit1
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb2() {
	cpu.operand &= ^bit2
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb3() {
	cpu.operand &= ^bit3
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb4() {
	cpu.operand &= ^bit4
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb5() {
	cpu.operand &= ^bit5
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb6() {
	cpu.operand &= ^bit6
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) rmb7() {
	cpu.operand &= ^bit7
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb0() {
	cpu.operand |= bit0
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb1() {
	cpu.operand |= bit1
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb2() {
	cpu.operand |= bit2
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb3() {
	cpu.operand |= bit3
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb4() {
	cpu.operand |= bit4
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb5() {
	cpu.operand |= bit5
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb6() {
	cpu.operand |= bit6
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) smb7() {
	cpu.operand |= bit7
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) brk() {
	pch := uint8((cpu.Registers.ProgramCounter >> 8) & 0xff)
	pcl := uint8(cpu.Registers.ProgramCounter & 0xff)
	cpu.pushStack(pch)
	cpu.pushStack(pcl)
	cpu.setStatusBit(BreakBit, true)
	cpu.pushStack(cpu.Registers.Status)
	cpu.setStatusBit(BreakBit, false)
	cpu.setStatusBit(InterruptDisableBit, true)
	cpu.setStatusBit(DecimalBit, false)
	cpu.Registers.ProgramCounter = uint16(cpu.read(vectorBRKH)) << 8
	cpu.Registers.ProgramCounter = cpu.Registers.ProgramCounter | uint16(cpu.read(vectorBRKL))
}

func (cpu *CPU) clc() {
	cpu.setStatusBit(CarryBit, false)
}

func (cpu *CPU) cld() {
	cpu.setStatusBit(DecimalBit, false)
}

func (cpu *CPU) cli() {
	cpu.setStatusBit(InterruptDisableBit, false)
}

func (cpu *CPU) clv() {
	cpu.setStatusBit(OverflowBit, false)
}

func (cpu *CPU) cmp() {
	res := cpu.Registers.Accumulator - cpu.operand
	cpu.setStatusBit(CarryBit, cpu.Registers.Accumulator >= cpu.operand)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == cpu.operand)
	cpu.setStatusBit(NegativeBit, res&0x80 != 0)
}

func (cpu *CPU) cpy() {
	res := cpu.Registers.Y - cpu.operand
	cpu.setStatusBit(CarryBit, cpu.Registers.Y >= cpu.operand)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Y == cpu.operand)
	cpu.setStatusBit(NegativeBit, res&0x80 != 0)
}

func (cpu *CPU) cpx() {
	res := cpu.Registers.X - cpu.operand
	cpu.setStatusBit(CarryBit, cpu.Registers.X >= cpu.operand)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == cpu.operand)
	cpu.setStatusBit(NegativeBit, res&0x80 != 0)
}

func (cpu *CPU) dec() {
	cpu.operand--
	cpu.setStatusBit(NegativeBit, cpu.operand&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

//DEC - Accumulator
func (cpu *CPU) deca() {
	cpu.Registers.Accumulator--
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) dex() {
	cpu.Registers.X--
	cpu.setStatusBit(NegativeBit, cpu.Registers.X&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == 0)
}

func (cpu *CPU) dey() {
	cpu.Registers.Y--
	cpu.setStatusBit(NegativeBit, cpu.Registers.Y&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Y == 0)
}

func (cpu *CPU) eor() {
	cpu.Registers.Accumulator = cpu.Registers.Accumulator ^ cpu.operand
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) inc() {
	cpu.operand++
	cpu.setStatusBit(NegativeBit, cpu.operand&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

//INC - Accumulator
func (cpu *CPU) inca() {
	cpu.Registers.Accumulator++
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) inx() {
	cpu.Registers.X++
	cpu.setStatusBit(NegativeBit, cpu.Registers.X&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == 0)
}

func (cpu *CPU) iny() {
	cpu.Registers.Y++
	cpu.setStatusBit(NegativeBit, cpu.Registers.Y&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Y == 0)
}

func (cpu *CPU) jmp() {
	cpu.operandAddress++
	pch := uint16(cpu.read(cpu.operandAddress)) << 8
	cpu.Registers.ProgramCounter = pch | uint16(cpu.operand)
}

func (cpu *CPU) jsr() {
	pch := uint8((cpu.Registers.ProgramCounter & 0xff) >> 8)
	pcl := uint8(cpu.Registers.ProgramCounter & 0xff)
	cpu.pushStack(pch)
	cpu.pushStack(pcl)
	cpu.operandAddress++
	cpu.Registers.ProgramCounter = (uint16(cpu.read(cpu.operandAddress)) << 8) | uint16(cpu.operand)
}

func (cpu *CPU) lda() {
	cpu.Registers.Accumulator = cpu.operand
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) ldx() {
	cpu.Registers.X = cpu.operand
	cpu.setStatusBit(NegativeBit, cpu.Registers.X&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == 0)
}

func (cpu *CPU) ldy() {
	cpu.Registers.Y = cpu.operand
	cpu.setStatusBit(NegativeBit, cpu.Registers.Y&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Y == 0)
}

func (cpu *CPU) lsr() {
	cpu.setStatusBit(CarryBit, cpu.operand&0x01 != 0)
	cpu.operand = cpu.operand >> 1
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.setStatusBit(NegativeBit, false)
	cpu.write(cpu.operandAddress, cpu.operand)
}

//LSR - Accumulator
func (cpu *CPU) lsra() {
	cpu.setStatusBit(CarryBit, cpu.Registers.Accumulator&0x01 != 0)
	cpu.Registers.Accumulator = cpu.Registers.Accumulator >> 1
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
	cpu.setStatusBit(NegativeBit, false)
}

func (cpu *CPU) ora() {
	cpu.Registers.Accumulator |= cpu.operand
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) pha() {
	cpu.pushStack(cpu.Registers.Accumulator)
}

func (cpu *CPU) php() {
	cpu.pushStack(cpu.Registers.Status)
}

func (cpu *CPU) phx() {
	cpu.pushStack(cpu.Registers.X)
}

func (cpu *CPU) phy() {
	cpu.pushStack(cpu.Registers.Y)
}

func (cpu *CPU) pla() {
	cpu.Registers.Accumulator = cpu.pullStack()
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) plp() {
	cpu.Registers.Status = cpu.pullStack()
}

func (cpu *CPU) plx() {
	cpu.Registers.X = cpu.pullStack()
	cpu.setStatusBit(NegativeBit, cpu.Registers.X&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == 0)
}

func (cpu *CPU) ply() {
	cpu.Registers.Y = cpu.pullStack()
	cpu.setStatusBit(NegativeBit, cpu.Registers.Y&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Y == 0)
}

func (cpu *CPU) rol() {
	carryIn := uint8(0)
	if cpu.testStatusBit(CarryBit) {
		carryIn++
	}
	cpu.setStatusBit(CarryBit, cpu.operand&0x80 != 0)
	cpu.operand = (cpu.operand << 1) | carryIn
	cpu.setStatusBit(NegativeBit, cpu.operand&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

//ROL - Accumulator
func (cpu *CPU) rola() {
	carryIn := uint8(0)
	if cpu.testStatusBit(CarryBit) {
		carryIn++
	}
	cpu.setStatusBit(CarryBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.Registers.Accumulator = (cpu.Registers.Accumulator << 1) | carryIn
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) ror() {
	carryIn := uint8(0)
	if cpu.testStatusBit(CarryBit) {
		carryIn = 0x80
	}
	cpu.setStatusBit(CarryBit, cpu.operand&0x01 == 0x01)
	cpu.operand = (cpu.operand >> 1) | carryIn
	cpu.setStatusBit(NegativeBit, cpu.operand&0x80 == 0x80)
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

//ROR - Accumulator
func (cpu *CPU) rora() {
	carryIn := uint8(0)
	if cpu.testStatusBit(CarryBit) {
		carryIn = 0x80
	}
	cpu.setStatusBit(CarryBit, cpu.Registers.Accumulator&0x01 == 0x01)
	cpu.Registers.Accumulator = (cpu.Registers.Accumulator >> 1) | carryIn
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 == 0x80)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) rti() {
	cpu.Registers.Status = cpu.pullStack()
	cpu.Registers.ProgramCounter = uint16(cpu.pullStack())
	cpu.Registers.ProgramCounter |= uint16(cpu.pullStack()) << 8
	if cpu.handlingNMI {
		cpu.handlingNMI = false
	}
}

func (cpu *CPU) rts() {
	cpu.Registers.ProgramCounter = uint16(cpu.pullStack())
	cpu.Registers.ProgramCounter |= uint16(cpu.pullStack()) << 8
	cpu.Registers.ProgramCounter++
}

func (cpu *CPU) sec() {
	cpu.setStatusBit(CarryBit, true)
}

func (cpu *CPU) sed() {
	cpu.setStatusBit(DecimalBit, true)
}

func (cpu *CPU) sei() {
	cpu.setStatusBit(InterruptDisableBit, true)
}

func (cpu *CPU) sta() {
	cpu.write(cpu.operandAddress, cpu.Registers.Accumulator)
}

func (cpu *CPU) stx() {
	cpu.write(cpu.operandAddress, cpu.Registers.X)
}

func (cpu *CPU) sty() {
	cpu.write(cpu.operandAddress, cpu.Registers.Y)
}

func (cpu *CPU) stz() {
	cpu.write(cpu.operandAddress, 0)
}

func (cpu *CPU) tax() {
	cpu.Registers.X = cpu.Registers.Accumulator
	cpu.setStatusBit(NegativeBit, cpu.Registers.X&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == 0)
}

func (cpu *CPU) tay() {
	cpu.Registers.Y = cpu.Registers.Accumulator
	cpu.setStatusBit(NegativeBit, cpu.Registers.Y&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Y == 0)
}

func (cpu *CPU) trb() {
	cpu.operand = ^cpu.Registers.Accumulator & cpu.operand
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) tsb() {
	cpu.operand = cpu.Registers.Accumulator & cpu.operand
	cpu.setStatusBit(ZeroBit, cpu.operand == 0)
	cpu.write(cpu.operandAddress, cpu.operand)
}

func (cpu *CPU) tsx() {
	cpu.Registers.X = cpu.Registers.StackPointer
	cpu.setStatusBit(NegativeBit, cpu.Registers.X&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.X == 0)
}

func (cpu *CPU) txa() {
	cpu.Registers.Accumulator = cpu.Registers.X
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) txs() {
	cpu.Registers.StackPointer = cpu.Registers.X
}

func (cpu *CPU) tya() {
	cpu.Registers.Accumulator = cpu.Registers.Y
	cpu.setStatusBit(NegativeBit, cpu.Registers.Accumulator&0x80 != 0)
	cpu.setStatusBit(ZeroBit, cpu.Registers.Accumulator == 0)
}

func (cpu *CPU) wai() {
	cpu.waiting = true
}

func (cpu *CPU) stp() {
	cpu.stopped = true
}

/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Instruction lookup table

NOTES: this lookup table was generated automatically. Each array element
corresponds to one CPU opcode. All illegal opcodes are mapped to NOP as per
WDC specifications.
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

type instruction struct {
	addressing func(*CPU)
	operation  func(*CPU)
}

var instructionTable = [256]instruction{
	{nil, (*CPU).brk},         //0
	{(*CPU).zpii, (*CPU).ora}, //1
	{nil, nil},                //2
	{nil, nil},                //3
	{(*CPU).zp, (*CPU).tsb},   //4
	{(*CPU).zp, (*CPU).ora},   //5
	{(*CPU).zp, (*CPU).asl},   //6
	{(*CPU).zp, (*CPU).rmb0},  //7
	{nil, (*CPU).php},         //8
	{(*CPU).imm, (*CPU).ora},  //9
	{nil, (*CPU).asla},        //10
	{nil, nil},                //11
	{(*CPU).abs, (*CPU).tsb},  //12
	{(*CPU).abs, (*CPU).ora},  //13
	{(*CPU).abs, (*CPU).asl},  //14
	{(*CPU).zp, (*CPU).bbr0},  //15
	{nil, (*CPU).bpl},         //16
	{(*CPU).zpiy, (*CPU).ora}, //17
	{(*CPU).zpi, (*CPU).ora},  //18
	{nil, nil},                //19
	{(*CPU).zp, (*CPU).trb},   //20
	{(*CPU).zpx, (*CPU).ora},  //21
	{(*CPU).zpx, (*CPU).asl},  //22
	{(*CPU).zp, (*CPU).rmb1},  //23
	{nil, (*CPU).clc},         //24
	{(*CPU).aiy, (*CPU).ora},  //25
	{nil, (*CPU).inca},        //26
	{nil, nil},                //27
	{(*CPU).abs, (*CPU).trb},  //28
	{(*CPU).aix, (*CPU).ora},  //29
	{(*CPU).aix, (*CPU).asl},  //30
	{(*CPU).zp, (*CPU).bbr1},  //31
	{(*CPU).abs, (*CPU).jsr},  //32
	{(*CPU).zpii, (*CPU).and}, //33
	{nil, nil},                //34
	{nil, nil},                //35
	{(*CPU).zp, (*CPU).bit},   //36
	{(*CPU).zp, (*CPU).and},   //37
	{(*CPU).zp, (*CPU).rol},   //38
	{(*CPU).zp, (*CPU).rmb2},  //39
	{nil, (*CPU).plp},         //40
	{(*CPU).imm, (*CPU).and},  //41
	{nil, (*CPU).rola},        //42
	{nil, nil},                //43
	{(*CPU).abs, (*CPU).bit},  //44
	{(*CPU).abs, (*CPU).and},  //45
	{(*CPU).abs, (*CPU).rol},  //46
	{(*CPU).zp, (*CPU).bbr2},  //47
	{nil, (*CPU).bmi},         //48
	{(*CPU).zpiy, (*CPU).and}, //49
	{(*CPU).zp, (*CPU).and},   //50
	{nil, nil},                //51
	{(*CPU).zpx, (*CPU).bit},  //52
	{(*CPU).zpx, (*CPU).and},  //53
	{(*CPU).zpx, (*CPU).rol},  //54
	{(*CPU).zp, (*CPU).rmb3},  //55
	{nil, (*CPU).sec},         //56
	{(*CPU).aiy, (*CPU).and},  //57
	{nil, (*CPU).deca},        //58
	{nil, nil},                //59
	{(*CPU).aix, (*CPU).bit},  //60
	{(*CPU).aix, (*CPU).and},  //61
	{(*CPU).aix, (*CPU).rol},  //62
	{(*CPU).zp, (*CPU).bbr3},  //63
	{nil, (*CPU).rti},         //64
	{(*CPU).zpii, (*CPU).eor}, //65
	{nil, nil},                //66
	{nil, nil},                //67
	{nil, nil},                //68
	{(*CPU).zp, (*CPU).eor},   //69
	{(*CPU).zp, (*CPU).lsr},   //70
	{(*CPU).zp, (*CPU).rmb4},  //71
	{nil, (*CPU).pha},         //72
	{(*CPU).imm, (*CPU).eor},  //73
	{nil, (*CPU).lsra},        //74
	{nil, nil},                //75
	{(*CPU).abs, (*CPU).jmp},  //76
	{(*CPU).abs, (*CPU).eor},  //77
	{(*CPU).abs, (*CPU).lsr},  //78
	{(*CPU).abs, (*CPU).bbr4}, //79
	{nil, (*CPU).bvc},         //80
	{(*CPU).zpiy, (*CPU).eor}, //81
	{(*CPU).zpi, (*CPU).eor},  //82
	{nil, nil},                //83
	{nil, nil},                //84
	{(*CPU).zpx, (*CPU).eor},  //85
	{(*CPU).zpx, (*CPU).lsr},  //86
	{(*CPU).zp, (*CPU).rmb5},  //87
	{nil, (*CPU).cli},         //88
	{(*CPU).aiy, (*CPU).eor},  //89
	{nil, (*CPU).phy},         //90
	{nil, nil},                //91
	{nil, nil},                //92
	{(*CPU).aix, (*CPU).eor},  //93
	{(*CPU).aix, (*CPU).lsr},  //94
	{(*CPU).zp, (*CPU).bbr5},  //95
	{nil, (*CPU).rts},         //96
	{(*CPU).zpii, (*CPU).adc}, //97
	{nil, nil},                //98
	{nil, nil},                //99
	{(*CPU).zp, (*CPU).stz},   //100
	{(*CPU).zp, (*CPU).adc},   //101
	{(*CPU).zp, (*CPU).ror},   //102
	{(*CPU).zp, (*CPU).rmb6},  //103
	{nil, (*CPU).pla},         //104
	{(*CPU).imm, (*CPU).adc},  //105
	{nil, (*CPU).rora},        //106
	{nil, nil},                //107
	{(*CPU).ai, nil},          //108 jmp -> absolute indirect
	{(*CPU).abs, (*CPU).adc},  //109
	{(*CPU).abs, (*CPU).ror},  //110
	{(*CPU).zp, (*CPU).bbr6},  //111
	{nil, (*CPU).bvs},         //112
	{(*CPU).zpiy, (*CPU).adc}, //113
	{(*CPU).zpi, (*CPU).adc},  //114
	{nil, nil},                //115
	{(*CPU).zpx, (*CPU).stz},  //116
	{(*CPU).zpx, (*CPU).adc},  //117
	{(*CPU).zpx, (*CPU).ror},  //118
	{(*CPU).zp, (*CPU).rmb7},  //119
	{nil, (*CPU).sei},         //120
	{(*CPU).aiy, (*CPU).adc},  //121
	{nil, (*CPU).ply},         //122
	{nil, nil},                //123
	{(*CPU).aii, (*CPU).jmp},  //124
	{(*CPU).aix, (*CPU).adc},  //125
	{(*CPU).aix, (*CPU).ror},  //126
	{(*CPU).zp, (*CPU).bbr7},  //127
	{nil, (*CPU).pcr},         //128
	{(*CPU).zpii, (*CPU).sta}, //129
	{nil, nil},                //130
	{nil, nil},                //131
	{(*CPU).zp, (*CPU).sty},   //132
	{(*CPU).zp, (*CPU).sta},   //133
	{(*CPU).zp, (*CPU).stx},   //134
	{(*CPU).zp, (*CPU).smb0},  //135
	{nil, (*CPU).dey},         //136
	{(*CPU).imm, (*CPU).bit},  //137
	{nil, (*CPU).txa},         //138
	{nil, nil},                //139
	{(*CPU).abs, (*CPU).sty},  //140
	{(*CPU).abs, (*CPU).sta},  //141
	{(*CPU).abs, (*CPU).stx},  //142
	{(*CPU).zp, (*CPU).bbs0},  //143
	{nil, (*CPU).bcc},         //144
	{(*CPU).zpiy, (*CPU).sta}, //145
	{(*CPU).zpi, (*CPU).sta},  //146
	{nil, nil},                //147
	{(*CPU).zpx, (*CPU).sty},  //148
	{(*CPU).zpx, (*CPU).sta},  //149
	{(*CPU).zpy, (*CPU).stx},  //150
	{(*CPU).zp, (*CPU).smb1},  //151
	{nil, (*CPU).tya},         //152
	{(*CPU).aiy, (*CPU).sta},  //153
	{nil, (*CPU).txs},         //154
	{nil, nil},                //155
	{(*CPU).abs, (*CPU).stz},  //156
	{(*CPU).aix, (*CPU).sta},  //157
	{(*CPU).aix, (*CPU).stz},  //158
	{(*CPU).zp, (*CPU).bbs1},  //159
	{(*CPU).imm, (*CPU).ldy},  //160
	{(*CPU).zpii, (*CPU).lda}, //161
	{(*CPU).imm, (*CPU).ldx},  //162
	{nil, nil},                //163
	{(*CPU).zp, (*CPU).ldy},   //164
	{(*CPU).zp, (*CPU).lda},   //165
	{(*CPU).zp, (*CPU).ldx},   //166
	{(*CPU).zp, (*CPU).smb2},  //167
	{nil, (*CPU).tay},         //168
	{(*CPU).imm, (*CPU).lda},  //169
	{nil, (*CPU).tax},         //170
	{nil, nil},                //171
	{nil, (*CPU).ldy},         //172
	{(*CPU).abs, (*CPU).lda},  //173
	{(*CPU).abs, (*CPU).ldx},  //174
	{(*CPU).zp, (*CPU).bbs2},  //175
	{nil, (*CPU).bcs},         //176
	{(*CPU).zpiy, (*CPU).lda}, //177
	{(*CPU).zpi, (*CPU).lda},  //178
	{nil, nil},                //179
	{(*CPU).zpx, (*CPU).ldy},  //180
	{(*CPU).zpx, (*CPU).lda},  //181
	{(*CPU).zpy, (*CPU).ldx},  //182
	{(*CPU).zp, (*CPU).smb3},  //183
	{nil, (*CPU).clv},         //184
	{(*CPU).aiy, (*CPU).lda},  //185
	{nil, (*CPU).tsx},         //186
	{nil, nil},                //187
	{(*CPU).aix, (*CPU).ldy},  //188
	{(*CPU).aix, (*CPU).lda},  //189
	{(*CPU).aiy, (*CPU).ldx},  //190
	{(*CPU).zp, (*CPU).bbs3},  //191
	{(*CPU).imm, (*CPU).cpy},  //192
	{(*CPU).zpii, (*CPU).cmp}, //193
	{nil, nil},                //194
	{nil, nil},                //195
	{(*CPU).zp, (*CPU).cpy},   //196
	{(*CPU).zp, (*CPU).cmp},   //197
	{(*CPU).zp, (*CPU).dec},   //198
	{(*CPU).zp, (*CPU).smb4},  //199
	{nil, (*CPU).iny},         //200
	{(*CPU).imm, (*CPU).cmp},  //201
	{nil, (*CPU).dex},         //202
	{nil, (*CPU).wai},         //203
	{(*CPU).abs, (*CPU).cpy},  //204
	{(*CPU).abs, (*CPU).cmp},  //205
	{(*CPU).abs, (*CPU).dec},  //206
	{(*CPU).zp, (*CPU).bbs4},  //207
	{nil, (*CPU).bne},         //208
	{(*CPU).zpiy, (*CPU).cmp}, //209
	{(*CPU).zpi, (*CPU).cmp},  //210
	{nil, nil},                //211
	{nil, nil},                //212
	{(*CPU).zpx, (*CPU).cmp},  //213
	{(*CPU).zpx, (*CPU).dec},  //214
	{(*CPU).zp, (*CPU).smb5},  //215
	{nil, (*CPU).cld},         //216
	{(*CPU).aiy, (*CPU).cmp},  //217
	{nil, (*CPU).phx},         //218
	{nil, (*CPU).stp},         //219
	{nil, nil},                //220
	{(*CPU).aix, (*CPU).cmp},  //221
	{(*CPU).aix, (*CPU).dec},  //222
	{(*CPU).zp, (*CPU).bbs5},  //223
	{(*CPU).imm, (*CPU).cpx},  //224
	{(*CPU).zpii, (*CPU).sbc}, //225
	{nil, nil},                //226
	{nil, nil},                //227
	{(*CPU).zp, (*CPU).cpx},   //228
	{(*CPU).zp, (*CPU).sbc},   //229
	{(*CPU).zp, (*CPU).inc},   //230
	{(*CPU).zp, (*CPU).smb6},  //231
	{nil, (*CPU).inx},         //232
	{(*CPU).imm, (*CPU).sbc},  //233
	{nil, nil},                //234
	{nil, nil},                //235
	{(*CPU).abs, (*CPU).cpx},  //236
	{(*CPU).abs, (*CPU).sbc},  //237
	{(*CPU).abs, (*CPU).inc},  //238
	{(*CPU).zp, (*CPU).bbs6},  //239
	{nil, (*CPU).beq},         //240
	{(*CPU).zpiy, (*CPU).sbc}, //241
	{(*CPU).zpi, (*CPU).sbc},  //242
	{nil, nil},                //243
	{nil, nil},                //244
	{(*CPU).zpx, (*CPU).sbc},  //245
	{(*CPU).zpx, (*CPU).inc},  //246
	{(*CPU).zp, (*CPU).smb7},  //247
	{nil, (*CPU).sed},         //248
	{(*CPU).aiy, (*CPU).sbc},  //249
	{nil, (*CPU).plx},         //250
	{nil, nil},                //251
	{nil, nil},                //252
	{(*CPU).aix, (*CPU).sbc},  //253
	{(*CPU).aix, (*CPU).inc},  //254
	{(*CPU).zp, (*CPU).bbs7},  //255
}
