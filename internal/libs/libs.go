package libs

import (
	"debug/elf"
	"fmt"
)

func Run(args []string) error {
	for _, file := range args {
		if err := run(file); err != nil {
			return err
		}
	}
	return nil

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
