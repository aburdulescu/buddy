package main

import (
	"flag"
	"fmt"
	"os"

	"bandr.me/p/buddy/internal/libs"
	"bandr.me/p/buddy/internal/syms"
)

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func mainErr() error {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: buddy COMMAND

Commands:
  libs  Print imported libraries
  syms  Print symbols
`)
	}
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		return fmt.Errorf("need command")
	}

	cmd := args[0]
	args = args[1:]

	switch flag.Arg(0) {
	case "libs":
		return libs.Run(args)
	case "syms":
		return syms.Run(args)
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}
