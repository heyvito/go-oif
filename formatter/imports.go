package formatter

import (
	"bytes"
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"
)

type ImportKind int

func (i ImportKind) String() string {
	switch i {
	case ImportKindLocal:
		return "Local"
	case ImportKindBuiltin:
		return "Builtin"
	case ImportKindThirdParty:
		return "ThirdParty"
	default:
		return "Unknown"
	}
}

const (
	ImportKindBuiltin ImportKind = iota
	ImportKindThirdParty
	ImportKindLocal
)

type Import struct {
	Import *dst.ImportSpec
	Kind   ImportKind
}

func (i Import) String() string {
	return fmt.Sprintf("%s: %s", i.Kind.String(), i.Import.Path.Value)
}

var stdImports = map[string]bool{}

func init() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}
	for _, p := range pkgs {
		stdImports[p.PkgPath] = true
	}
}

func isStdImport(path string) bool {
	_, ok := stdImports[path]
	return ok
}

func detectImportKind(projectPath string, importPath string) ImportKind {
	switch {
	case isStdImport(importPath):
		return ImportKindBuiltin
	case strings.HasPrefix(importPath, projectPath):
		return ImportKindLocal
	default:
		return ImportKindThirdParty
	}
}

func FormatImports(projectPath string, src string) string {
	f, err := decorator.Parse(src)
	if err != nil {
		panic(err)
	}

	var firstImport, lastImport *dst.GenDecl
	var imports []Import
	topDecorators := dst.NodeDecs{
		Before: -1,
	}

	for _, rd := range f.Decls {
		d, ok := rd.(*dst.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			// Not an import declaration, so we're done.
			// Imports are always first.
			break
		}

		if firstImport == nil {
			firstImport = d
		}
		lastImport = d

		rdd := rd.Decorations()
		topDecorators.End.Append(rdd.End...)
		topDecorators.Start.Append(rdd.Start...)
		if topDecorators.Before == -1 {
			topDecorators.Before = rd.Decorations().Before
		}

		for _, s := range d.Specs {
			if spec, ok := s.(*dst.ImportSpec); ok {
				imports = append(imports, Import{
					Import: spec,
					Kind:   detectImportKind(projectPath, importPath(spec)),
				})
			}
		}
	}

	var builtinImports, thirdPartyImports, localImports []Import

	for _, i := range imports {
		switch i.Kind {
		case ImportKindBuiltin:
			builtinImports = append(builtinImports, i)
		case ImportKindThirdParty:
			thirdPartyImports = append(thirdPartyImports, i)
		case ImportKindLocal:
			localImports = append(localImports, i)
		}
	}

	if firstImport == nil || lastImport == nil {
		return src
	}
	specs := generateImports(builtinImports, thirdPartyImports, localImports)
	multi := len(specs) > 1

	importsDecl := dst.GenDecl{
		Tok:    token.IMPORT,
		Lparen: multi,
		Rparen: multi,
		Specs:  specs,
		Decs: dst.GenDeclDecorations{
			NodeDecs: topDecorators,
			Tok:      nil,
			Lparen:   nil,
		},
	}

	appending := true
	var newFile []dst.Decl
	for _, i := range f.Decls {
		if i == firstImport {
			appending = false
			newFile = append(newFile, &importsDecl)
		}
		if i == lastImport {
			appending = true
		}

		if i == lastImport || i == firstImport {
			continue
		}

		if appending {
			newFile = append(newFile, i)
		}
	}

	f.Decls = newFile

	buf := bytes.Buffer{}
	rest := decorator.NewRestorer()
	if err = rest.Fprint(&buf, f); err != nil {
		panic(err)
	}
	return buf.String()
}

func generateImports(builtin, third, local []Import) []dst.Spec {
	var res []dst.Spec

	items := [][]Import{builtin, third, local}
	for _, imps := range items {
		for _, x := range imps {
			i := x.Import
			i.Decorations().After = dst.NewLine
			i.Decorations().Before = dst.NewLine
			res = append(res, x.Import)
		}
		if len(res) > 0 {
			res[len(res)-1].Decorations().After = dst.EmptyLine
		}
	}
	return res
}

func importPath(s dst.Spec) string {
	t, err := strconv.Unquote(s.(*dst.ImportSpec).Path.Value)
	if err == nil {
		return t
	}
	return ""
}
