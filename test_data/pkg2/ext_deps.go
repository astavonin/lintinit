package pkg2

import (
	"fmt"

	"github.com/astavonin/lintinit/test_data/pkg2/sub"

	"github.com/astavonin/lintinit/test_data/pkg1"
	sh "github.com/astavonin/lintinit/test_data/shared"
	//. "github.com/astavonin/lintinit/test_data/pkg1"
)

var globalData = "some data"
var BooVal = false

func init() {
	globalData = pkg1.GlobalData

	// nolint
	pkg1.GlobalDataInt++

	if pkg1.GlobalDataInt > 0 || pkg1.GlobalStruct.BooVal < 0 {
		pkg1.Foo()
	}

	if pkg1.GlobalDeepStruct.BooStruct.BooVal != 42 && !BooVal {
		fmt.Println(sh.SomeFunc(), sub.SubExport())
	}
}

func foo() {
	fmt.Println(globalData)
}
