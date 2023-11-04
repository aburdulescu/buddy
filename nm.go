package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ianlancetaylor/demangle"
)

var (
	noDemangle  = flag.Bool("no-demangle", false, "Do not demangle C++/Rust names")
	symbolTable = flag.String("t", "symtab", "Symbol table to read: symtab or dynsym")
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		if err := run(path); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}

func run(path string) error {
	file, err := elf.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var symbols []elf.Symbol

	switch *symbolTable {
	case "symtab":
		symbols, err = file.Symbols()
	case "dynsym":
		symbols, err = file.DynamicSymbols()
	default:
		return fmt.Errorf("unknown value for -t: %s", *symbolTable)
	}

	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Binding\tType\tVisibility\tLibrary\tVersion\tName")
	fmt.Fprintln(w, "-------\t----\t----------\t-------\t-------\t----")

	for _, sym := range symbols {
		name := sym.Name
		if !*noDemangle {
			name = demangle.Filter(name)
		}
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			elf.ST_BIND(sym.Info),
			elf.ST_TYPE(sym.Info),
			elf.ST_VISIBILITY(sym.Other),
			sym.Library,
			sym.Version,
			name,
		)
	}

	return nil
}
