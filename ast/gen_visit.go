//go:build ignore

package main

import (
	"bytes"
	"cmp"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"slices"
)

// Generates visit.go

type NodeType int

const (
	NodeTypeStruct NodeType = iota
	NodeTypeSlice
)

type VisitableNodeType struct {
	Type     NodeType
	Name     string
	Children []Child
}

type Child struct {
	FieldName string
	Optional  bool
}

func newChild(fieldName string, optional bool) Child {
	return Child{FieldName: fieldName, Optional: optional}
}

func main() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "./ast", func(info fs.FileInfo) bool {
		return info.Name() != "visit.go"
	}, parser.ParseComments)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var nodes []VisitableNodeType
	for _, file := range pkgs["ast"].Files {
		nodes = append(nodes, findVisitableNodes(file)...)
	}

	slices.SortFunc(nodes, func(a, b VisitableNodeType) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fmt.Println(nodes)

	var (
		visitorMethods     []*ast.Field
		noopVisitorMethods []ast.Decl
		visitMethods       []ast.Decl
	)
	for _, node := range nodes {
		visitorMethods = append(visitorMethods, &ast.Field{
			Names: []*ast.Ident{{Name: "Visit" + node.Name}},
			Type: &ast.FuncType{
				Params: newFieldList("n", &ast.StarExpr{X: ast.NewIdent(node.Name)}),
			},
		})

		noopVisitorMethods = append(noopVisitorMethods, &ast.FuncDecl{
			Recv: newFieldList("nv", &ast.StarExpr{X: ast.NewIdent("NoopVisitor")}),
			Name: ast.NewIdent("Visit" + node.Name),
			Type: &ast.FuncType{
				Params: newFieldList("n", &ast.StarExpr{X: ast.NewIdent(node.Name)}),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun:  newSelectorExpr(ast.NewIdent("n"), "VisitChildrenWith"),
						Args: []ast.Expr{newSelectorExpr(ast.NewIdent("nv"), "V")},
					}},
				},
			},
		})

		recv := newFieldList("n", &ast.StarExpr{X: ast.NewIdent(node.Name)})
		params := newFieldList("v", ast.NewIdent("Visitor"))
		visitChildrenBlock := &ast.BlockStmt{}
		switch node.Type {
		case NodeTypeStruct:
			for _, child := range node.Children {
				callExpr := &ast.ExprStmt{X: &ast.CallExpr{
					Fun: newSelectorExpr(
						newSelectorExpr(ast.NewIdent("n"), child.FieldName),
						"VisitWith",
					),
					Args: []ast.Expr{ast.NewIdent("v")},
				}}
				if child.Optional {
					visitChildrenBlock.List = append(visitChildrenBlock.List, &ast.IfStmt{
						Cond: &ast.BinaryExpr{
							X:  newSelectorExpr(ast.NewIdent("n"), child.FieldName),
							Op: token.NEQ,
							Y:  ast.NewIdent("nil"),
						},
						Body: &ast.BlockStmt{List: []ast.Stmt{callExpr}},
					})
				} else {
					visitChildrenBlock.List = append(visitChildrenBlock.List, callExpr)
				}
			}
		case NodeTypeSlice:
			// for i := 0; i < len(*n); i++ {
			//     (*n)[i].VisitWith(v)
			// }
			visitChildrenBlock.List = append(visitChildrenBlock.List, &ast.ForStmt{
				Init: &ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("i")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{&ast.BasicLit{Value: "0"}},
				},
				Cond: &ast.BinaryExpr{
					X:  ast.NewIdent("i"),
					Op: token.LSS,
					Y: &ast.CallExpr{
						Fun:  ast.NewIdent("len"),
						Args: []ast.Expr{&ast.StarExpr{X: ast.NewIdent("n")}},
					},
				},
				Post: &ast.IncDecStmt{
					X:   ast.NewIdent("i"),
					Tok: token.INC,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{X: &ast.CallExpr{
							Fun: newSelectorExpr(&ast.IndexExpr{
								X:     &ast.StarExpr{X: ast.NewIdent("n")},
								Index: ast.NewIdent("i"),
							}, "VisitWith"),
							Args: []ast.Expr{ast.NewIdent("v")},
						}},
					},
				},
			})
		}
		visitMethods = append(visitMethods, &ast.FuncDecl{
			Recv: recv,
			Name: ast.NewIdent("VisitWith"),
			Type: &ast.FuncType{Params: params},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{X: &ast.CallExpr{
						Fun:  newSelectorExpr(ast.NewIdent("v"), "Visit"+node.Name),
						Args: []ast.Expr{ast.NewIdent("n")},
					}},
				},
			},
		}, &ast.FuncDecl{
			Recv: recv,
			Name: ast.NewIdent("VisitChildrenWith"),
			Type: &ast.FuncType{Params: params},
			Body: visitChildrenBlock,
		})
	}

	genPkg := &ast.File{
		Name: ast.NewIdent("ast"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Visitor"),
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{List: visitorMethods},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("NoopVisitor"),
						Type: &ast.StructType{
							Fields: newFieldList("V", ast.NewIdent("Visitor")),
						},
					},
				},
			},
		},
	}

	genPkg.Decls = append(genPkg.Decls, noopVisitorMethods...)
	genPkg.Decls = append(genPkg.Decls, visitMethods...)

	s := bytes.NewBuffer([]byte("// Code generated by gen_visit.go; DO NOT EDIT.\n"))
	format.Node(s, fset, genPkg)

	os.WriteFile("ast/visit.go", s.Bytes(), 0644)

	fmt.Println(pkgs)
}

func findVisitableNodes(f *ast.File) (types []VisitableNodeType) {
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			switch typeSpec.Name.Name {
			case "ScopeContext", "Id":
				continue
			}

			switch t := typeSpec.Type.(type) {
			case *ast.StructType:
				types = append(types, VisitableNodeType{
					Type:     NodeTypeStruct,
					Name:     typeSpec.Name.Name,
					Children: findStructChildren(t.Fields.List),
				})
			case *ast.ArrayType:
				types = append(types, VisitableNodeType{
					Type: NodeTypeSlice,
					Name: typeSpec.Name.Name,
				})
			}
		}
	}
	return types
}

func findStructChildren(fields []*ast.Field) (children []Child) {
	for _, field := range fields {
		optional := field.Tag != nil && field.Tag.Value == "`optional:\"true\"`"

		if len(field.Names) != 0 {
			fmt.Println(field.Names[0].Name)
		}

		switch fieldType := field.Type.(type) {
		case *ast.Ident:
			if len(field.Names) == 0 {
				children = append(children, newChild(fieldType.Name, optional))
				continue
			}

			switch fieldType.Name {
			case "Idx", "any", "bool", "int", "ScopeContext", "string", "PropertyKind", "float64":
			default:
				fmt.Println(fieldType.Name)
				children = append(children, newChild(field.Names[0].Name, optional))
			}
		case *ast.StarExpr:
			if ident, ok := fieldType.X.(*ast.Ident); ok && ident.Name == "string" {
				continue
			}
			children = append(children, newChild(field.Names[0].Name, optional))
		}
	}
	return children
}

func newFieldList(name string, t ast.Expr) *ast.FieldList {
	return &ast.FieldList{
		List: []*ast.Field{{
			Names: []*ast.Ident{ast.NewIdent(name)},
			Type:  t,
		}},
	}
}

func newSelectorExpr(x ast.Expr, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: x, Sel: ast.NewIdent(sel)}
}
