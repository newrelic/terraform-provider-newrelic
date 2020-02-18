package crudfuncinfo

import (
	"go/ast"
	"reflect"

	"github.com/bflad/tfproviderlint/helper/astutils"
	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "crudfuncinfo",
	Doc:  "find github.com/hashicorp/terraform-plugin-sdk/helper/schema CreateFunc, ReadFunc, UpdateFunc, and DeleteFunc declarations for later passes",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run:        run,
	ResultType: reflect.TypeOf([]*schema.CRUDFuncInfo{}),
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
	}
	var result []*schema.CRUDFuncInfo

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcType := astutils.FuncTypeFromNode(n)

		if funcType == nil {
			return
		}

		if !astutils.IsFieldListTypePackageType(funcType.Params, 0, pass.TypesInfo, schema.PackagePath, schema.TypeNameResourceData) {
			return
		}

		if !astutils.IsFieldListType(funcType.Params, 1, astutils.IsFunctionParameterTypeInterface) {
			return
		}

		if !astutils.IsFieldListType(funcType.Results, 0, astutils.IsFunctionParameterTypeError) {
			return
		}

		result = append(result, schema.NewCRUDFuncInfo(n, pass.TypesInfo))
	})

	return result, nil
}
