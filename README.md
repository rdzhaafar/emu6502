# Emu6502

Emu6502 is a W65C02S CPU emulation library written in Go. It features a full implementation of the W65C02S ISA, as well as a simple command-line debugger.

## Project goals

**This project's primary goal is to be an educational tool for anyone wanting to learn 6502 assembly programming and architecture**

6502 is great, because it's easy to understand in comparison to modern CPU's, yet sophisticated enough to be useful. While it's hard to get your hands on the original model, variants of it are still being produced today. WDC's W65C02S is one of them. It's cheap enough to be disposable, making it a go-to for DIY projects and learning.

This library implements the W65C02S instruction set, which adds some new instructions on top of the original 6502 ISA. It was designed to be as simple to understand as possible. If your primary focus is speed, then you should use one of many other 6502 emulators instead.

_Note that while this library emulates some electronic behaviour of the W65C02 like sending interrupts, it does not emulate ALL electronic behaviour such as sending clock signals._

## How to get the library

[Go](https://golang.org/) needs to be installed on your computer in order to run this command. Get the library with
```shell
go get github.com/rdzhaafar/emu6502
```

## How to use the library

For instructions on how to use the library, read the [core library guide](core).

## How to use the debugger

For instructions on how to use the debugger, read the [command line debugger guide](debugger).

to start the debugger.

## Licence
[MIT](LICENCE)