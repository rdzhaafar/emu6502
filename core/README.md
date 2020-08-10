# Core library guide

Import core library into your project

```go
import "github.com/rdzhaafar/emu6502/core"
```

Initialize the system bus and write the program to it. For this example, we're going to write 

```asm
    LDA #88 ;load the accumulator with "88"
```

which, translated to 6502 machine code is

```asm
    A9 88
```

```go
bus := core.NewBasicBus()
bus.Write(0x0000, 0xA9)
bus.Write(0x0001, 0x88)
```

Then, initialize the registers and lastly the CPU itself.

```go
registers := core.NewCPURegisters()
cpu := core.NewCPU(bus, registers)
```

You can then execute the program by using the `Execute` function.

```go
cpu.Execute()
```

You can send an interrupt to the cpu using `cpu.Interrupt`. To send a non-maskable interrupt, use `cpu.NMInterrupt`. To reset the cpu use `cpu.Reset`.
