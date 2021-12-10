// Copyright (c) 2021, Geert JM Vanderkelen

package xdoc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/geertjanvdk/xkit/xlog"
)

type DocInfo struct {
	PackageDoc  string
	PackageName string
	Constants   Values
	Variables   Values
	Functions   Functions
	Structs     Structs
}

func isExported(name string) bool {
	c, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(c)
}

func New(toDoc string) *DocInfo {
	fileSet := token.NewFileSet()

	dir, err := parser.ParseDir(fileSet, toDoc, nil, parser.ParseComments)
	if err != nil {
		xlog.Error(err)
	}

	docInfo := &DocInfo{
		Constants: make(Values),
		Variables: make(Values),
		Functions: make(Functions),
		Structs:   make(Structs),
	}

	for packageName, a := range dir {
		if strings.HasSuffix(packageName, "_test") {
			continue
		}
		docInfo.PackageName = packageName

		// first pass: get all struct types, constants, and variables
		for filename, file := range a.Files {
			if strings.HasSuffix(filename, "_test.go") {
				continue
			}

			if strings.Contains(filename, path.Join(docInfo.PackageName, "doc.go")) && file.Doc != nil {
				docInfo.PackageDoc = file.Doc.Text()
			}

			for _, decl := range file.Decls {
				d, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				var doc string
				if d.Doc != nil {
					doc = d.Doc.Text()
				}

				switch d.Tok.String() {
				case "type":
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.TypeSpec:
							if !isExported(s.Name.String()) {
								continue
							}

							st, ok := s.Type.(*ast.StructType)
							if !ok {
								continue
							}

							n := newStruct(s.Name.String())

							if st.Fields != nil {
								for _, f := range st.Fields.List {
									if len(f.Names) == 0 {
										// embedded
										n.Fields = append(n.Fields, whatIsExpr(f.Type))
										continue
									}

									n.Fields = append(n.Fields, fmt.Sprintf("%s %s %s",
										f.Names[0],
										whatIsExpr(f.Type),
										whatIsExpr(f.Tag),
									))
								}
							}

							n.Doc = doc
							docInfo.Structs[packageName+"."+s.Name.String()] = n
						}
					}
				case "const", "var":
					for _, spec := range d.Specs {
						switch s := spec.(type) {
						case *ast.ValueSpec:
							if !isExported(s.Names[0].Name) {
								continue
							}
							c := newValue(s.Names[0].String(), whatIsExpr(s.Values[0]))
							c.IsConstant = d.Tok.String() == "const"
							c.Doc = doc
							if c.IsConstant {
								docInfo.Constants[packageName+"."+c.Name] = c
							} else {
								docInfo.Variables[packageName+"."+c.Name] = c
							}
						}
					}
				}
			}
		}

		// second pass: get functions and method; separating Constructors
		for filename, file := range a.Files {
			if strings.HasSuffix(filename, "_test.go") {
				continue
			}

			if filename == path.Join(packageName, "doc.go") && file.Doc != nil {
				docInfo.PackageDoc = ""
				for _, d := range file.Doc.List {
					txt := strings.TrimSpace(d.Text)
					txt = strings.Replace(txt, "/*", "", 1)
					txt = strings.Replace(txt, "*/", "", 1)
					txt = strings.TrimSpace(txt)
					docInfo.PackageDoc += txt
				}
			}

			for _, decl := range file.Decls {
				var registered bool

				d, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}

				if !d.Name.IsExported() {
					continue
				}

				var doc string
				if d.Doc != nil {
					doc = d.Doc.Text()
				}

				if d.Recv != nil {
					// this is a method (function which has a receiver)
					n := strings.TrimPrefix(whatIsExpr(d.Recv.List[0].Type), "*")
					if !isExported(n) {
						continue
					}
					if t, have := docInfo.Structs[packageName+"."+n]; have {
						t.Methods[d.Name.String()] = &Function{
							Name:  d.Name.String(),
							type_: d.Type,
							recv:  d.Recv,
							Doc:   doc,
						}
						registered = true
					}
				} else {
					// check if this function is a constructor, or it returns a struct define in package
					if d.Type.Results != nil {
						for _, r := range d.Type.Results.List {
							n := packageName + "." + strings.TrimPrefix(whatIsExpr(r.Type), "*")
							if t, have := docInfo.Structs[n]; have {
								t.Constructors[d.Name.String()] = &Function{
									Name:  d.Name.String(),
									type_: d.Type,
									recv:  d.Recv,
									Doc:   doc,
								}
								registered = true
								break
							}
						}
					}

					// non-constructor function
					if !registered {
						docInfo.Functions[packageName+"."+d.Name.String()] = &Function{
							Name:  d.Name.String(),
							type_: d.Type,
							recv:  d.Recv,
							Doc:   doc,
						}
					}
				}
			}
		}
	}

	return docInfo
}

func whatIsExpr(x ast.Expr) string {
	if x == nil {
		return ""
	}
	switch n := x.(type) {
	case *ast.ArrayType:
		return "[]" + whatIsExpr(n.Elt)
	case *ast.Ellipsis:
		return "..." + whatIsExpr(n.Elt)
	case *ast.Ident:
		if ast.IsExported(n.Name) {
			return n.Name
		}
		return n.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		return whatIsExpr(n.X) + "." + n.Sel.String()
	case *ast.StarExpr:
		return "*" + whatIsExpr(n.X)
	case *ast.FuncType:
		return "func" + funcSignature("", n, nil)
	case *ast.BasicLit:
		if n == nil {
			return ""
		}
		return n.Value
	default:
		return fmt.Sprintf("%T", x)
	}
}

func argumentList(fields *ast.FieldList) []string {
	if fields == nil || fields.List == nil {
		return nil
	}

	params := make([]string, len(fields.List))
	for i, p := range fields.List {
		expr := whatIsExpr(p.Type)
		if p.Names != nil {
			params[i] = p.Names[0].Name + " " + expr
		} else {
			params[i] = expr
		}
	}

	return params
}

func funcSignature(fnName string, ft *ast.FuncType, recv *ast.FieldList) string {
	params := argumentList(ft.Params)
	results := argumentList(ft.Results)

	var result string
	switch len(results) {
	case 0:
		result = ""
	case 1:
		result = " " + results[0]
	default:
		result = " (" + strings.Join(results, ", ") + ")"
	}

	var receiver string
	if recv != nil {
		receiver = "(" + argumentList(recv)[0] + ") "
	}

	return fmt.Sprintf("%s%s(%s)%s", receiver, fnName, strings.Join(params, ", "), result)
}
