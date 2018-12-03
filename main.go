package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func processFolder(root string) ([]LintError, error) {
	var (
		lintErrs []LintError
		err      error
	)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		l := NewLinter(path, nil)
		pkgErrs, err := l.Parse()
		if err != nil {
			return err
		}
		lintErrs = append(lintErrs, pkgErrs...)
		return nil
	})

	return lintErrs, err
}

func printErrors(errors []LintError) {
	for _, err := range errors {
		fmt.Printf("%s:%d:%d: global resource %s access from from package init() call\n",
			err.Line.File, err.Line.Line, err.Line.Column, err.Ident)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of lintinit:\n")
	fmt.Fprintf(os.Stderr, "\tlintinit [directory] # runs on package in current or [directory] recursively\n")
	flag.PrintDefaults()
}

func main() {
	//flag.Usage = usage
	//flag.Parse()

	root := "."

	le, err := processFolder(root)
	if err != nil {
		log.Fatal(err)
	}

	printErrors(le)
}
