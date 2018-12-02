package main

import (
	"fmt"
	"os"
)

func main() {
	l := NewLinter(os.Args[1], nil)
	fmt.Println(l.Parse())
}

// https://github.com/golang/go/blob/master/src/cmd/vet/unused.go
// https://github.com/golang/example/tree/master/gotypes
// https://stackoverflow.com/questions/32532335/usage-of-go-parser-across-packages
// https://arslan.io/2017/09/14/the-ultimate-guide-to-writing-a-go-tool/
// https://stackoverflow.com/questions/50524607/go-lang-func-parameter-type
//
