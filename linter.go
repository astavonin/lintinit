package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strings"
)

type Package struct {
	Name string
	Path []string
}

func ParsePackage(rawPath, shortName string) *Package {
	fullName := strings.Trim(rawPath, "\"")

	parsed := strings.Split(fullName, "/")
	importName := shortName
	if importName == "" {
		importName = parsed[len(parsed)-1]
	}

	return &Package{importName, parsed}
}

func (p *Package) String() string {
	return "Package(" + p.Name + ", " + strings.Join(p.Path, "/") + ")"
}

// IsChild returns true if current package is child for pkgRoot
func (p *Package) IsChildFor(pkgRoot *Package) bool {
	if len(pkgRoot.Path) >= len(p.Path) {
		return false
	}

	for i, val := range pkgRoot.Path {
		if p.Path[i] != val {
			return false
		}
	}
	return true
}

func (p *Package) IsInternal() bool {
	return len(p.Path) > 1
}

type Ident struct {
	Path     []string
	position token.Position
}

func (i *Ident) String() string {
	return "Ident(" + strings.Join(i.Path, ".") + ", " + i.position.String() + ")"
}

func NewIdent(names []string, position token.Position) *Ident {
	return &Ident{names, position}
}

func (i *Ident) PkgName() string {
	if len(i.Path) <= 1 {
		return "."
	}
	return i.Path[0]
}

func (i *Ident) Name() string {
	return i.Path[len(i.Path)-1]
}

func (i *Ident) FullName() string {
	return strings.Join(i.Path, ".")
}

type LineInfo struct {
	File   string
	Line   int
	Column int
}

func (li LineInfo) String() string {
	return fmt.Sprintf("LineInfo(%s:%d:%d)", li.File, li.Line, li.Column)
}

func NewLineInfo(position token.Position) LineInfo {
	return LineInfo{File: position.Filename, Line: position.Line, Column: position.Column}
}

type LintError struct {
	Line  LineInfo
	Ident string
}

func (le LintError) String() string {
	return fmt.Sprintf("LintError(Ident=%s, Line=%s)", le.Ident, le.Line)
}

type Linter interface {
	Parse() ([]LintError, error)
	Path() string
	PkgName() string
}

type linter struct {
	dir   string
	defs  []string
	files []string

	root    *Package
	imports map[string]*Package
	idents  []*Ident
}

func (l *linter) Path() string {
	panic("implement me")
}

func (l *linter) PkgName() string {
	panic("implement me")
}

func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

func (l *linter) buildFilesList() error {
	ctx := build.Default
	ctx.BuildTags = l.defs

	pkg, err := ctx.ImportDir(l.dir, 0)
	if err != nil {
		return fmt.Errorf("cannot process directory %s: %s", l.dir, err)
	}
	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.SFiles...)

	l.files = prefixDirectory(l.dir, files)
	l.root = ParsePackage(pkg.ImportPath, "")
	return nil
}

func (l *linter) parse() (*types.Package, error) {

	fset := token.NewFileSet()
	var astFiles []*ast.File

	for _, fname := range l.files {
		f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("parsing error: %s, %s", fname, err)
		}
		astFiles = append(astFiles, f)

		//fmt.Println("AAAA:", f.Name.String())
		//ast.Print(fset, f)

		ast.Inspect(f, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.ImportSpec:
				pkg := parseImport(t)
				l.imports[pkg.Name] = pkg
				return false
			case *ast.FuncDecl:
				if t.Name.Name == "init" {
					l.idents = append(l.idents, processInit(fset, t.Body)...)
					return false
				}
			}
			return true
		})
	}
	if len(astFiles) == 0 {
		return nil, fmt.Errorf("%s: no Go files found", l.dir)
	}

	config := types.Config{
		IgnoreFuncBodies: false,
		Importer:         importer.For("source", nil),
		FakeImportC:      true,
	}

	return config.Check(l.dir, fset, astFiles, nil)
}

func parseImport(spec *ast.ImportSpec) *Package {
	importName := ""
	if spec.Name != nil {
		importName = spec.Name.Name
	}

	return ParsePackage(spec.Path.Value, importName)
}

func collectFromSelectors(sel *ast.SelectorExpr) (res []string, pos token.Pos) {
	switch t := sel.X.(type) {
	case *ast.SelectorExpr:
		r, p := collectFromSelectors(t)
		res = append(r, sel.Sel.Name)
		pos = p
	case *ast.Ident:
		res = append(res, t.Name, sel.Sel.Name)
		pos = sel.Sel.NamePos
	}
	return
}

func collectFromIdents(ident *ast.Ident) []string {
	return []string{ident.Name}
}

func processInit(fset *token.FileSet, decl *ast.BlockStmt) []*Ident {
	//ast.Print(fset, decl)
	var acc []*Ident
	ast.Inspect(decl, func(n ast.Node) bool {
		deeper := true
		switch t := n.(type) {
		case *ast.SelectorExpr:
			res, pos := collectFromSelectors(t)
			acc = append(acc, NewIdent(res, fset.Position(pos)))
			deeper = false
		case *ast.Ident:
			acc = append(acc, NewIdent(collectFromIdents(t), fset.Position(t.NamePos)))
			deeper = false
		}
		return deeper
	})
	return acc
}

func (l *linter) Parse() ([]LintError, error) {

	err := l.buildFilesList()

	if err != nil {
		return nil, err
	}

	_, err = l.parse()
	if err != nil {
		return nil, err
	}
	//
	//fmt.Println("Parse:", l.root)
	//
	//fmt.Println("Imports:", l.imports)
	//fmt.Println("Identifiers:", l.idents)
	//
	//for _, Ident := range l.idents {
	//	fmt.Println("ID:", Ident.PkgName())
	//}

	//fmt.Println("Pkg imports:", pkg.Imports())
	//fmt.Println("Lookup:", pkg.Scope().Lookup("BooVal"))

	var lintErr []LintError
	for _, ident := range l.idents {
		if ident.PkgName() == "." {
			// local scope
			continue
		}
		pkg, ok := l.imports[ident.PkgName()]
		if !ok {
			log.Fatal(fmt.Sprintf("unknown Ident %s", ident))
		}
		if pkg.IsInternal() && !pkg.IsChildFor(l.root) {
			lintErr = append(lintErr, LintError{
				NewLineInfo(ident.position),
				ident.FullName()})
		}
	}

	return lintErr, nil
}

func NewLinter(dir string, defs []string) Linter {
	return &linter{
		dir:     dir,
		defs:    defs,
		imports: map[string]*Package{},
	}
}
