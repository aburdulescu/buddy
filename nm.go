package main

import (
	"debug/elf"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/ianlancetaylor/demangle"
)

var (
	demangleMode = flag.String("demangle", "short", "How to demangle C++/Rust names: none, short, full")
	symbolTable  = flag.String("table", "symtab", "Symbol table to read: symtab or dynsym")
	outFile      = flag.String("o", "", "Path to output file")
	writeJSON    = flag.Bool("json", false, "Write output as JSON")
)

func main() {
	flag.Parse()
	if err := run(flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(path string) error {
	switch *demangleMode {
	case "none", "short", "full":
	default:
		return fmt.Errorf("unknown value for -demangle: %s", *demangleMode)
	}

	file, err := elf.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var rawSymbols []elf.Symbol

	switch *symbolTable {
	case "symtab":
		rawSymbols, err = file.Symbols()
	case "dynsym":
		rawSymbols, err = file.DynamicSymbols()
	default:
		return fmt.Errorf("unknown value for -table: %s", *symbolTable)
	}

	if err != nil {
		return err
	}

	var out io.Writer = os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	}

	var nameFilter func(name string) string

	switch *demangleMode {
	case "none":
		nameFilter = func(name string) string { return name }
	case "short":
		opts := []demangle.Option{
			demangle.NoParams,
			demangle.NoTemplateParams,
			demangle.NoEnclosingParams,
			demangle.NoRust,
		}
		nameFilter = func(name string) string {
			return demangle.Filter(name, opts...)
		}
	case "full":
		nameFilter = func(name string) string {
			return demangle.Filter(name)
		}
	}

	symbols := make([]Symbol, len(rawSymbols))
	for i, sym := range rawSymbols {
		symbols[i] = Symbol{
			Name:       nameFilter(sym.Name),
			Binding:    elf.ST_BIND(sym.Info).String(),
			Type:       elf.ST_TYPE(sym.Info).String(),
			Visibility: elf.ST_VISIBILITY(sym.Other).String(),
			Library:    sym.Library,
			Version:    sym.Version,
			Value:      sym.Value,
			Size:       sym.Size,
		}
	}

	if *writeJSON {
		if err := json.NewEncoder(out).Encode(symbols); err != nil {
			return err
		}
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 2, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Binding\tType\tVisibility\tLibrary\tVersion\tValue\tSize\tName")
	fmt.Fprintln(w, "-------\t----\t----------\t-------\t-------\t-----\t----\t----")

	for _, sym := range symbols {
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s\n",
			sym.Binding,
			sym.Type,
			sym.Visibility,
			sym.Library,
			sym.Version,
			sym.Value,
			sym.Size,
			sym.Name,
		)
	}

	return nil
}

type Symbol struct {
	Name       string
	Binding    string
	Type       string
	Visibility string
	Library    string
	Version    string
	Value      uint64
	Size       uint64
}
