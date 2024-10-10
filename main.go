package main

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	celast "github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/parser"
	"log"
	"os"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func selects(mef cel.MacroExprFactory, base celast.Expr, fields ...celast.Expr) celast.Expr {
	for _, field := range fields {
		base = mef.NewSelect(base, field.AsIdent())
	}
	return base
}
func foldl[T any](slice []T,  combine func(T, T) T) T {
	result := slice[0]
	for _, v := range slice[1:] {
		result = combine(result, v)
	}
	return result
}

//func simpleBuilder(name string,)

func main() {
	expr := os.Args[1]
	prsr, err := parser.NewParser(
		parser.Macros(
			// Usage: `self.index(x, z, b)`
			// This does a nil-safe traversal of
			cel.ReceiverVarArgMacro("index",
				func(mef cel.MacroExprFactory, base celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
					if len(args) == 0 {
						return nil, mef.NewError(base.ID(), "index requires at least 1 arg")
					}
					checks := []celast.Expr{}
					for i  := range args{

						next := mef.NewCall(operators.Has, selects(mef, base, args[0:i+1]...))
						checks = append(checks, next)
						logExpr(fmt.Sprintf("iter %d", i), next)
					}
					check := foldl(checks, func(l, r celast.Expr) celast.Expr {
						return mef.NewCall(operators.LogicalAnd, l, r)
					})

					final := mef.NewCall(operators.Conditional, check, selects(mef, base, args...), mef.NewLiteral(types.NullValue))
					return final, nil
				}),
			// Usage: `oneof(self.x, self.y, self.z)`
			// This checks that 0 or 1 of these fields is set, mirroring Protobuf one-of checking logic.
			cel.GlobalVarArgMacro("oneof",
				func(mef cel.MacroExprFactory, base celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
					if len(args) < 2 {
						return nil, mef.NewError(base.ID(), "oneof requires at least 2 arg")
					}
					checks := []celast.Expr{}
					for _, arg  := range args{

						has := mef.NewCall(operators.Has, arg)
					logExpr("sum", arg)
						check := mef.NewCall(operators.Conditional, has, mef.NewLiteral(types.Int(1)), mef.NewLiteral(types.Int(0)))
						checks = append(checks, check)
					}
					sum := foldl(checks, func(l, r celast.Expr) celast.Expr {
						return mef.NewCall(operators.Add, l, r)
					})

					final := mef.NewCall(operators.LessEquals, sum, mef.NewLiteral(types.Int(1)))
					return final, nil
				}),
			// Usage: `default(self.x, 'DEF')`.
			// This returns self.x if its set, else 'DEF'
			cel.GlobalMacro("default", 2,
				func(mef cel.MacroExprFactory, iterRange celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
					has := mef.NewCall(operators.Has, args[0])
					return mef.NewCall(operators.Conditional, has, args[0], args[1]), nil
				}),
		),
	)
	fatal(err)
	p, iss := prsr.Parse(common.NewTextSource(expr))
	if len(iss.GetErrors()) > 0 {
		fatal(err)
	}
	c, _ := celast.ToProto(p)
	s, _ := (&jsonpb.Marshaler{Indent: "  "}).MarshalToString(c)
	log.Println(s)
	out, err := parser.Unparse(p.Expr(), p.SourceInfo())
	fatal(err)
	fmt.Println(out)
}

func logExpr(base string, expr celast.Expr) {
	out, err := parser.Unparse(expr, nil)
	fatal(err)
	log.Println(base, out)
}