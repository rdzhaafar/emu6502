package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rdzhaafar/emu6502/core"
)

type shellOptions struct {
	loadBinaryFile bool
	binaryFileName string
}

type shellCommand struct {
	command string
	args    []string
}

type interactiveShell struct {
	options *shellOptions
	cpu     *core.CPU
}

func newShellCommand(input []byte) *shellCommand {
	input = bytes.Trim(input, " ")
	str := strings.Split(string(input), " ")
	if len(str) == 0 {
		return nil
	}
	cmd := shellCommand{}
	cmd.command = strings.ToLower(str[0])
	if len(str) > 1 {
		args := str[1:]
		for i, arg := range args {
			args[i] = strings.ToLower(arg)
		}
		cmd.args = args
	} else {
		cmd.args = nil
	}
	return &cmd
}

func (shell *interactiveShell) handleCommand(cmd *shellCommand) bool {
	switch cmd.command {
	case "print":
		shell.printCmd(cmd)
		return false
	case "set":
		shell.setCmd(cmd)
		return false
	case "step":
		shell.stepCmd(cmd)
		return false
	case "load":
		shell.loadCmd(cmd)
		return false
	case "exit":
		return true
	case "help":
		shell.helpCmd(cmd)
		return false
	default:
		shell.invalidCmd(cmd.command)
		return false
	}
}

func (shell *interactiveShell) printError(cmd string, err string) {
	fmt.Printf("Error: %s: %s\n", cmd, err)
}

func (shell *interactiveShell) printInfo(cmd string, info string) {
	fmt.Printf("Info: %s: %s\n", cmd, info)
}

func (shell *interactiveShell) invalidCmd(cmd string) {
	err := fmt.Sprintf("%s is not a valid command.", cmd)
	shell.printError(cmd, err)
}

func (shell *interactiveShell) invalidArgs(cmd string, args []string) {
	err := fmt.Sprintf("%s are not valid arguments for this command", args)
	shell.printError(cmd, err)
}

func (shell *interactiveShell) printCmd(cmd *shellCommand) {
	args := cmd.args
	switch l := len(args); l {
	case 1:
		switch args[0] {
		case "registers":
			shell.printRegisters()
		default:
			shell.invalidArgs(cmd.command, args)
		}

	case 2:
		switch args[0] {
		case "bus":
			addr, err := strconv.ParseUint(args[1], 16, 16)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.printBus(uint16(addr))

		default:
			shell.invalidArgs(cmd.command, args)
		}

	case 3:
		addrl, err := strconv.ParseUint(args[1], 16, 16)
		if err != nil {
			shell.invalidArgs(cmd.command, args)
			break
		}
		addrh, err := strconv.ParseUint(args[2], 16, 16)
		if err != nil {
			shell.invalidArgs(cmd.command, args)
			break
		}
		if addrh < addrl {
			shell.invalidArgs(cmd.command, args)
			break
		}
		shell.printBusRange(uint16(addrl), uint16(addrh))

	default:
		shell.invalidArgs(cmd.command, args)
	}
}

func (shell *interactiveShell) helpCmd(cmd *shellCommand) {
	switch len(cmd.args) {
	case 0:
		fmt.Println("Interactive shell commands:")
		fmt.Println("1. step -> steps the cpu through one instruction")
		fmt.Println("2. print -> prints the contents of cpu registers or the bus")
		fmt.Println("\t\"print registers\" prints the contents of cpu registers")
		fmt.Println("\t\"print bus X\" prints the contents of the bus at address X")
		fmt.Println("\t\"print bus X Y\" prints the contents of the bus from address X to Y")
		fmt.Println("3. set -> sets the registers/bus to the specified value")
		fmt.Println("\t\"set (a, p, x, y, pc, sp) X\" sets the the specified register to X")
		fmt.Println("\t\"set bus X Y\" sets the bus at address X to Y")
		fmt.Println("4. load -> loads a binary file to the bus for debugging")
		fmt.Println("\t\"load X\" loads file X (where X is either an absolute path, or a relative path)")
		fmt.Println("5. help -> prints this help message")
		fmt.Println("6. exit -> exits the interactive shell")
	default:
		shell.invalidArgs(cmd.command, cmd.args)
	}
}

func (shell *interactiveShell) printRegisters() {
	fmt.Printf("A:  %02X\n", shell.cpu.Registers.Accumulator)
	fmt.Printf("X:  %02X\n", shell.cpu.Registers.X)
	fmt.Printf("Y:  %02X\n", shell.cpu.Registers.Y)
	fmt.Printf("P:  %08b\n", shell.cpu.Registers.Status)
	fmt.Printf("SP: (01)%02X\n", shell.cpu.Registers.StackPointer)
	fmt.Printf("PC: %04X\n", shell.cpu.Registers.ProgramCounter)
}

func (shell *interactiveShell) printBusRange(addrl uint16, addrh uint16) {
	for i := addrl; i <= addrh; i++ {
		shell.printBus(i)
	}
}

func (shell *interactiveShell) printBus(addr uint16) {
	val := shell.cpu.Bus.Read(uint16(addr))
	fmt.Printf("%04X: %02X\n", addr, val)
}

func (shell *interactiveShell) setCmd(cmd *shellCommand) {
	args := cmd.args
	switch l := len(args); l {
	case 2:
		switch args[0] {

		case "a":
			val, err := strconv.ParseUint(args[1], 16, 8)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Registers.Accumulator = uint8(val)
			shell.printInfo(cmd.command, fmt.Sprintf("Set accumulator to %02X", val))

		case "x":
			val, err := strconv.ParseUint(args[1], 16, 8)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Registers.X = uint8(val)
			shell.printInfo(cmd.command, fmt.Sprintf("Set x register to %02X", val))

		case "y":
			val, err := strconv.ParseUint(args[1], 16, 8)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Registers.Y = uint8(val)
			shell.printInfo(cmd.command, fmt.Sprintf("Set y register to %02X", val))

		case "p":
			val, err := strconv.ParseUint(args[1], 16, 8)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Registers.Status = uint8(val)
			shell.printInfo(cmd.command, fmt.Sprintf("Set status register to %02X", val))

		case "pc":
			val, err := strconv.ParseUint(args[1], 16, 16)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Registers.ProgramCounter = uint16(val)
			shell.printInfo(cmd.command, fmt.Sprintf("Set program counter to %02X", val))

		case "sp":
			val, err := strconv.ParseUint(args[1], 16, 8)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Registers.StackPointer = uint8(val)
			shell.printInfo(cmd.command, fmt.Sprintf("Set stack pointer to %02X", val))

		default:
			shell.invalidArgs(cmd.command, args)
		}

	case 3:
		switch args[0] {

		case "bus":
			addr, err := strconv.ParseUint(args[1], 16, 16)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			val, err := strconv.ParseUint(args[2], 16, 8)
			if err != nil {
				shell.invalidArgs(cmd.command, args)
				break
			}
			shell.cpu.Bus.Write(uint16(addr), uint8(val))
			shell.printInfo(cmd.command, fmt.Sprintf("Set bus address %04X to %02X", addr, val))
		default:
			shell.invalidArgs(cmd.command, args)
		}
	default:
		shell.invalidArgs(cmd.command, args)
	}
}

func (shell *interactiveShell) loadCmd(cmd *shellCommand) {
	args := cmd.args
	switch l := len(args); l {
	case 1:
		err := shell.loadFile(args[0])
		if err != nil {
			shell.printError(cmd.command, fmt.Sprintf("Could not load file %v.\n", args[0]))
			break
		}
	default:
		shell.invalidArgs(cmd.command, args)
	}
}

func (shell *interactiveShell) loadFile(filename string) error {
	file, err := os.Open(filename)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return err
	}
	bytes := make([]byte, core.MaxBusSize)
	read, err := file.Read(bytes)
	if err != nil {
		return err
	}
	shell.cpu.Bus = core.NewBasicBus()
	for i := 0; i <= read; i++ {
		shell.cpu.Bus.Write(uint16(i), bytes[i])
	}
	shell.printInfo("load", fmt.Sprintf("Read %v bytes from file %v.\n", read, filename))
	return nil
}

func (shell *interactiveShell) stepCmd(cmd *shellCommand) {
	args := cmd.args
	switch l := len(args); l {
	case 0:
		shell.cpu.Execute()
		shell.printInfo(cmd.command, "Executed 1 cpu cycle")
	case 1:
		steps, err := strconv.ParseUint(args[2], 10, 32)
		if err != nil {
			shell.invalidArgs(cmd.command, args)
			break
		}
		if steps <= 0 {
			shell.invalidArgs(cmd.command, args)
			break
		}
		for i := uint64(0); i <= steps; i++ {
			shell.cpu.Execute()
		}
	default:
		shell.invalidArgs(cmd.command, args)
	}
}

func newInteractiveShell(opt *shellOptions) (*interactiveShell, error) {
	bus := core.NewBasicBus()
	shell := interactiveShell{}
	if opt.loadBinaryFile {
		err := shell.loadFile(opt.binaryFileName)
		if err != nil {
			return nil, err
		}
	}
	registers := core.NewCPURegisters()
	cpu := core.NewCPU(bus, registers)
	shell.cpu = cpu
	shell.options = opt
	return &shell, nil
}

func (shell *interactiveShell) run() error {
	fmt.Printf("Interactive 65c02 debugger shell.\n")
	reader := bufio.NewReader(os.Stdin)
	exit := false
	for !exit {
		fmt.Printf(">")
		line, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		cmd := newShellCommand(line)
		if cmd != nil {
			exit = shell.handleCommand(cmd)
		}
	}
	return nil
}
