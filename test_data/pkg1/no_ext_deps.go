package pkg1

import "fmt"

type Boo struct {
	BooVal int
}

type DeepBoo struct {
	BooStruct Boo
}

var GlobalData = "some data"
var GlobalDataInt = 42
var GlobalStruct = Boo{42}
var GlobalDeepStruct = DeepBoo{Boo{42}}

func init() {
	GlobalData = "new data"
}

func Foo() {
	fmt.Println(GlobalData, GlobalDataInt, GlobalStruct, GlobalDeepStruct)
}
