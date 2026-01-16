package linker

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type PackageInfo struct {
	Name       string
	ImportPath string
	Exports    []ExportedFunc
}

type ExportedFunc struct {
	Name string
	Doc  string
	Args []string // Simple representation for PoC
	Ret  []string
}

// Inspect loads and analyzes a Go package
func Inspect(importPath string, dir string) (*PackageInfo, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Dir: dir,
	}

	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load package: %v", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package load errors")
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no package found")
	}

	pkg := pkgs[0]
	info := &PackageInfo{
		Name:       pkg.Name,
		ImportPath: pkg.PkgPath,
	}

	for _, syntax := range pkg.Syntax {
		for _, decl := range syntax.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !ast.IsExported(fn.Name.Name) {
				continue
			}

			// Capture documentation
			doc := fn.Doc.Text()

			// Simple Type Analysis (Improvements needed for complex types)
			var args []string
			if fn.Type.Params != nil {
				for _, param := range fn.Type.Params.List {
					typeName := types.ExprString(param.Type)
					for range param.Names {
						args = append(args, typeName)
					}
					// Handle unnamed parameters (rare in valid code but possible in signatures)
					if len(param.Names) == 0 {
						args = append(args, typeName)
					}
				}
			}

			info.Exports = append(info.Exports, ExportedFunc{
				Name: fn.Name.Name,
				Doc:  strings.TrimSpace(doc),
				Args: args,
			})
		}
	}

	return info, nil
}
