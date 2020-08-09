# Command line debugger guide

## Installation/Starting up

`cd` into the debugger source directory. The location will be different depending on your OS and account name, but generally it should be in
```shell
cd $USERHOME/go/src/github.com/rdzhaafar/emu6502/debugger
```
where `$USERHOME` is your user home directory. Run

```shell
go run .
```

to launch an interactive debugging shell. You can also run


```shell
go run . --file a.out
```

in order to load an assembled 6502 program to RAM. In order to assemble W65C02 binaries I recommend using [vasm](http://sun.hasenbraten.de/vasm/) assembler with `--wdc02` and `--Fbin` options.

>NOTE: If you want to compile and install the debugger to your computer permanently, read the [go install documentation](https://golang.org/cmd/go/). 

## Using the interactive shell

There are five major commands available in the shell.

### print

Print command can print the contents of the registers or the system bus. Run `print bus A A9` to print contents of the system bus from address **000A** to **00A9**. Run `print registers` to print contents of the CPU registers.

### set

Set command sets the contents of system bus or cpu registers to the specified value. Run `set pc 0` to set the program counter to **0**. Run `set bus A 1` to set the device at bus address **A** to **1**.

### step

Step command executes one cpu instruction.

### load

Load command loads a binary file to the system bus. Run `load a.out` to load file **a.out** to the system bus. Files can be specified either as a relative path or an absolute path.

### exit

To exit the shell, run `exit`.