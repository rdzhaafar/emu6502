# Emu6502

Emu6502 is a W65C02S CPU emulator written in Go.

## Project goals

This project's primary goal is to be an educational tool for anyone wanting to learn 6502 assembly programming and architecture

## Features

- Complete implementation of the WDC65C02S ISA
- Interrupt support
- Can be used as a library or as a standalone emulator

_Note that while this library emulates some electronic behaviour of the WDC65C02S like sending interrupts, it does not emulate ALL electronic behaviour such as sending clock signals._

## How to get the library

[Go](https://golang.org/) needs to be installed on your computer in order to run this command. Get the library with
```shell
go get github.com/rdzhaafar/emu6502
```

## How to use the library

For instructions on how to use the library, read the [core library manual](core).

## How to use the debugger

For instructions on how to use the debugger, read the [command line debugger manual](debugger).

## Licence
[MIT](LICENCE)
