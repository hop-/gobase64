package main

import (
	"gobase64/coder"
	"fmt"
	"io"
	"os"
)

func printHelp() {
	fmt.Println("base64 [option]")
	fmt.Println("  Options")
	fmt.Println("  -d: decrypting mode")
	fmt.Println("  -e: encrypting mode")
	fmt.Println("  -h: print this message")
}

func main() {
	var isEncrypt bool
	argsLen := len(os.Args[1:])
	switch argsLen {
	case 0:
		isEncrypt = true
	case 1:
		opt := os.Args[1]
		if opt == "-d" {
			isEncrypt = false
			break
		} else if opt == "-e" {
			isEncrypt = true
			break
		} else if opt == "-h" {
			printHelp()
			return
		}
		printHelp()
		os.Exit(1)
	default:
		printHelp()
		os.Exit(1)
	}

	var base64HandlerFunc func(io.Reader, io.Writer) error

	if isEncrypt {
		base64HandlerFunc = coder.EncodeStream
	} else {
		base64HandlerFunc = coder.DecodeStream
	}

	base64HandlerFunc(os.Stdin, os.Stdout)
}
