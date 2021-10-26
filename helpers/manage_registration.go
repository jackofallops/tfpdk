package helpers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type operation string

const (
	Register   operation = "register"
	Unregister operation = "unregister"
)

type resourceType string

const (
	TypeResource   resourceType = "resource"
	TypeDatasource resourceType = "datasource"
)

type astKey struct {
	KeyValuePos int
	KeyKind     token.Token
	KeyValue    string
	ValFun      string
}

var resourceNameToAdd string

var fSet *token.FileSet

// UpdateRegistration modifies the `registration.go` file to add (or remove?) items from the relevant blocks when
// a user adds (or removes?) a reource or datasource via `tfpdk resource` or `tfpdk datasource`
//
// servicePackagePath = The path from the route of the provider to the service package to which the item is to be registered
// resourceName = the model name of a typed resource, or the snakeCase name of the item
// resource = string value for the item type, one of TypeResource or TypeDatasource
// op = one of Register or Unregister // TODO - removing is a future concern so not yet implemented
// isTyped = true if the resources uses the TypedSDK
func UpdateRegistration(servicePackagePath string, resourceName string, resource resourceType, _ operation, isTyped bool) error {
	resourceNameToAdd = resourceName
	fSet = token.NewFileSet()
	regFilePath := fmt.Sprintf("%s/registration.go", strings.TrimSuffix(servicePackagePath, "/"))
	regFile, err := parser.ParseFile(fSet, regFilePath, nil, 0)
	if err != nil {
		return err
	}

	nodeName := normaliseNodeName(resource, isTyped)
	ast.Inspect(regFile, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if ok {
			if fn.Name.Name == nodeName {
				ast.Inspect(fn.Body.List[0], func(r ast.Node) bool {
					ret, ok := r.(*ast.ReturnStmt)
					if ok {
						ast.Inspect(ret.Results[0], func(n ast.Node) bool {
							if out, ok := n.(*ast.CompositeLit); ok {
								astutil.Apply(out, AppendResourceToRegistrationBlock(), nil)
							}
							return true
						})
					}
					return true
				})
			}
		}
		return true
	})

	outBuf := new(bytes.Buffer)
	if err := format.Node(outBuf, fSet, regFile); err != nil {
		return err
	}

	if err := os.WriteFile(regFilePath, outBuf.Bytes(), 0755); err != nil {
		return err
	}

	return nil
}

func normaliseNodeName(input resourceType, isTyped bool) string {
	switch input {
	case TypeDatasource:
		if isTyped {
			return "DataSources"
		} else {
			return "SupportedDataSources"
		}
	default:
		if isTyped {
			return "Resources"
		}
	}
	return "SupportedResources"
}

func newUnTypedASTReturnEntry(key string, value string, pos int) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key: &ast.BasicLit{
			ValuePos: token.Pos(pos),
			Kind:     token.STRING,
			Value:    fmt.Sprintf("%q", strings.Trim(key, "\"")),
		},
		Value: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: value,
			},
		},
	}
}

func UpdateRegistrationByNode(servicePackagePath string, resourceName string, resource resourceType, _ operation, isTyped bool) error {
	fSet := token.NewFileSet()
	regFilePath := fmt.Sprintf("%s/registration.go", strings.TrimSuffix(servicePackagePath, "/"))
	regFile, err := parser.ParseFile(fSet, regFilePath, nil, 0)
	if err != nil {
		return err
	}

	nodeName := normaliseNodeName(resource, isTyped)
	newKeyValue := TerraformResourceName("azurerm", resourceName)
	n := astutil.Apply(regFile, func(c *astutil.Cursor) bool {
		if d, ok := c.Parent().(*ast.FuncDecl); ok && d.Name.Name == nodeName && c.Name() == "Body" {
			node := c.Node()
			log.Printf("%+v", node)
			if _, ok := node.(*ast.BlockStmt); ok {
				elts := node.(*ast.BlockStmt).List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts
				log.Printf("%+v", elts)
				elts = append(elts, newUnTypedASTReturnEntry(newKeyValue, "testEntry", 0))
				node.(*ast.BlockStmt).List[0].(*ast.ReturnStmt).Results[0].(*ast.CompositeLit).Elts = elts
				c.Replace(node)
			}
		}
		return true
	}, nil)

	outBuf := new(bytes.Buffer)
	if err := format.Node(outBuf, fSet, n); err != nil {
		return err
	}

	if err := os.WriteFile(regFilePath, outBuf.Bytes(), 0755); err != nil {
		return err
	}

	return nil
}

func AppendResourceToRegistrationBlock() astutil.ApplyFunc {
	return func(c *astutil.Cursor) bool {
		m := c.Node()
		if t, ok := m.(*ast.KeyValueExpr); ok {
			alreadyPresent := false
			snakeName := TerraformResourceName("azurerm", resourceNameToAdd)
			if strings.Trim(t.Key.(*ast.BasicLit).Value, "\"") == snakeName {
				alreadyPresent = true
			}
			p := c.Parent().(*ast.CompositeLit)
			if len(p.Elts)-1 == c.Index() && !alreadyPresent {
				c.InsertAfter(newUnTypedASTReturnEntry(snakeName, "testEntry", int(m.(*ast.KeyValueExpr).Value.(*ast.CallExpr).Rparen)+4))
			}
			return false
		}
		return true
	}
}
