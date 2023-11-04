package main

import (
	"debug/elf"
	"fmt"
	"os"
)

func main() {
	for _, path := range os.Args[1:] {
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
	libs, err := file.ImportedLibraries()
	if err != nil {
		return err
	}
	for _, lib := range libs {
		fmt.Println(lib)
	}
	return nil
}
