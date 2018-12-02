package sub

import "fmt"

func boo() {
	fmt.Println("foo")
}

func SubExport() int {
	return 42
}
