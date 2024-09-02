/*
このファイルを更新したあとは、`go build -buildmode=plugin ./donotuseas.go`を実行し、donotuseas.soを生成してください

生成したdonotuseas.soは、.gitignoreで無視されるので注意です。

"plugin was built with a different version"というエラーが出た場合は、以下を確認してください。
- go.modで指定しているgoのバージョンと、`golangci-lint --version`で表示されるgoのバージョンが一致していること
- go.modで指定しているgolang.org/x/toolsのバージョンと、`golangci-lint version --debug | grep "golang.org/x/tools"`で表示されるgolang.org/x/toolsのバージョンが一致していること
*/
package donotuseas

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "donotuseas",
	Doc:  "check `arg` is not used as `param` in function arguments",
	Run:  run,
}

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

var argument string  // -arg flag
var parameter string // -param flag

func init() {
	Analyzer.Flags.StringVar(&argument, "arg", argument, "name of the argument type to restrict in function call")
	Analyzer.Flags.StringVar(&parameter, "param", parameter, "name of the parameter type to check")
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true // 関数呼び出し以外は無視する
			}

			var id *ast.Ident
			switch fun := callExpr.Fun.(type) {
			case *ast.Ident: // 関数呼び出しの場合
				id = fun
			case *ast.SelectorExpr: // メソッド呼び出しの場合
				id = fun.Sel
			}

			funcObj := pass.TypesInfo.ObjectOf(id)
			funcDecl, ok := funcObj.(*types.Func) // types.Funcにキャストすることで、関数シグネチャが取得できる
			if !ok {
				return true
			}

			// 関数シグネチャを取得
			funcSignature := funcDecl.Type().(*types.Signature)
			params := funcSignature.Params()

			// 引数をチェックする
			for i, arg := range callExpr.Args {
				if i >= params.Len() {
					return true
				}

				argType := pass.TypesInfo.TypeOf(arg)               // 実引数の型
				paramType, ok := params.At(i).Type().(*types.Named) // 仮引数の型 ※Named型にキャストすることで型名が取得できる
				if !ok {
					return true
				}

				// [TODO] 設定で、仮引数（型A）に対して実引数（型B）を渡している箇所を検出できるようにしたい
				fmt.Printf("argType: %s, paramType: %s\n", argType.String(), paramType.String())
				if strings.HasSuffix(argType.String(), argument) && strings.HasSuffix(paramType.String(), parameter) {
					pass.Reportf(arg.Pos(),
						"do not use %+s as %+s in function call",
						argType.String(), paramType.String())
				}
			}
			return true
		})
	}
	return nil, nil
}
