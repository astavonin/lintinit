package main

import (
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	_, b, _, _  = runtime.Caller(0)
	testData, _ = filepath.Abs(filepath.Join(filepath.Dir(b), "test_data"))

	noInit        = filepath.Join(testData, "pkg0")
	noExtDepsGo   = filepath.Join(testData, "pkg1")
	withExtDepsGo = filepath.Join(testData, "pkg2")
)

func Test_Package(t *testing.T) {
	Convey("Package should detect relations", t, func() {
		pkgRoot := ParsePackage("github.com/astavonin/lintinit/test_data/pkg2", "")
		pkgChild := ParsePackage("github.com/astavonin/lintinit/test_data/pkg2/sub", "")
		pkgExternal := ParsePackage("fmt", "")
		pkgInternal := ParsePackage("github.com/astavonin/lintinit/test_data/pkg1", "")

		So(pkgChild.IsInternal(), ShouldBeTrue)
		So(pkgExternal.IsInternal(), ShouldBeFalse)
		So(pkgInternal.IsInternal(), ShouldBeTrue)

		So(pkgRoot.IsChildFor(pkgChild), ShouldBeFalse)
		So(pkgChild.IsChildFor(pkgChild), ShouldBeFalse)
		So(pkgChild.IsChildFor(pkgRoot), ShouldBeTrue)

		So(pkgExternal.IsChildFor(pkgRoot), ShouldBeFalse)
		So(pkgInternal.IsChildFor(pkgRoot), ShouldBeFalse)
	})
}

func Test_Idents(t *testing.T) {
	Convey("Package should detect ", t, func() {
		ident := NewIdent([]string{"sub", "child_name"}, nil, token.Position{}, false)

		So(ident.FullName(), ShouldEqual, "sub.child_name")
		So(ident.Name(), ShouldEqual, "child_name")
		So(ident.PkgName(), ShouldEqual, "sub")
	})
}

func Test_LinterBaseCases(t *testing.T) {
	Convey("Linter should ignore no-init packages", t, func() {
		l := NewLinter(noInit, nil)

		So(l, ShouldNotBeNil)

		lintErrs, err := l.Parse()
		So(err, ShouldBeNil)

		So(lintErrs, ShouldBeEmpty)
	})

	Convey("Linter should ignore init with no deps package", t, func() {
		l := NewLinter(noExtDepsGo, nil)

		So(l, ShouldNotBeNil)

		lintErrs, err := l.Parse()
		So(err, ShouldBeNil)

		So(lintErrs, ShouldBeEmpty)
	})

	Convey("Linter should report init with deps packages", t, func() {
		l := NewLinter(withExtDepsGo, nil)

		So(l, ShouldNotBeNil)

		lintErrs, err := l.Parse()
		So(err, ShouldBeNil)

		So(lintErrs, ShouldHaveLength, 10)
	})
}
