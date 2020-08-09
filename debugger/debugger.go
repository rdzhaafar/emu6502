package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	exitSuccess          int = 0
	exitStartupError     int = 1
	exitStatusUnexpected int = 2
)

func parseCliArguments() *shellOptions {
	opt := shellOptions{}
	var binaryFileName string
	opt.loadBinaryFile = false
	flag.StringVar(&binaryFileName, "file", "", "65c02 executable to debug")
	flag.Parse()
	if binaryFileName != "" {
		opt.loadBinaryFile = true
	}
	opt.binaryFileName = binaryFileName
	return &opt
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Got runtime panic.\nTrace: %v\n", r)
			os.Exit(exitStatusUnexpected)
		}
	}()
	opt := parseCliArguments()
	shell, err := newInteractiveShell(opt)
	if err != nil {
		fmt.Printf("Failed to start the debugger.\nTrace: \n%v\n", err)
		os.Exit(exitStartupError)
	}
	err = shell.run()
	if err != nil {
		fmt.Printf("Got unexpected error.\nTrace: %v\n", err)
		os.Exit(exitStatusUnexpected)
	}
	os.Exit(exitSuccess)
}
