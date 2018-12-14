package pkg2

import (
	"fmt"

	"github.com/astavonin/lintinit/test_data/pkg1"
	"github.com/astavonin/lintinit/test_data/pkg2/sub"
	sh "github.com/astavonin/lintinit/test_data/shared"
	//. "github.com/astavonin/lintinit/test_data/pkg1"
)

type ExtDepsType struct {
	Val int
}

var globalData = "some data"
var BooVal = false
var ShVar sh.SomeType

func init() {

	globalData = pkg1.GlobalData

	t := ExtDepsType{1}

	ShVar.Val1 = t.Val
	ShVar.Val1 = 42
	ShVar.DoCall()

	pkg1.GlobalDataInt++ //nolint

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
