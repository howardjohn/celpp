package macros

import (
	"github.com/google/cel-go/cel"
	celast "github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

var All = []cel.Macro{
	Default,
	Oneof,
	Index,
}

var (
	// Default provides the default() custom macro.
	// Usage: `default(self.x, 'DEF')`.
	// This returns self.x if its set, else 'DEF'
	Default = cel.GlobalMacro("default", 2,
		func(mef cel.MacroExprFactory, iterRange celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
			has := mef.NewCall(operators.Has, args[0])
			return mef.NewCall(operators.Conditional, has, args[0], args[1]), nil
		})

	// Oneof provides the oneof() custom macro.
	// Usage: `oneof(self.x, self.y, self.z)`
	// This checks that 0 or 1 of these fields is set, mirroring Protobuf one-of checking logic.
	Oneof = cel.GlobalVarArgMacro("oneof",
		func(mef cel.MacroExprFactory, base celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
			if len(args) < 2 {
				return nil, mef.NewError(base.ID(), "oneof requires at least 2 arg")
			}
			checks := []celast.Expr{}
			for _, arg := range args {

				has := mef.NewCall(operators.Has, arg)
				check := mef.NewCall(operators.Conditional, has, mef.NewLiteral(types.Int(1)), mef.NewLiteral(types.Int(0)))
				checks = append(checks, check)
			}
			sum := foldl(checks, func(l, r celast.Expr) celast.Expr {
				return mef.NewCall(operators.Add, l, r)
			})

			final := mef.NewCall(operators.LessEquals, sum, mef.NewLiteral(types.Int(1)))
			return final, nil
		})

	// Index provides the index() custom macro.
	// Usage: `self.index(x, z, b)`
	// This does a nil-safe traversal of
	Index = cel.ReceiverVarArgMacro("index",
		func(mef cel.MacroExprFactory, base celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
			if len(args) == 0 {
				return nil, mef.NewError(base.ID(), "index requires at least 1 arg")
			}
			checks := []celast.Expr{}
			for i := range args {

				next := mef.NewCall(operators.Has, selects(mef, base, args[0:i+1]...))
				checks = append(checks, next)
			}
			check := foldl(checks, func(l, r celast.Expr) celast.Expr {
				return mef.NewCall(operators.LogicalAnd, l, r)
			})

			final := mef.NewCall(operators.Conditional, check, selects(mef, base, args...), mef.NewLiteral(types.NullValue))
			return final, nil
		})
)

func selects(mef cel.MacroExprFactory, base celast.Expr, fields ...celast.Expr) celast.Expr {
	for _, field := range fields {
		base = mef.NewSelect(base, field.AsIdent())
	}
	return base
}
func foldl[T any](slice []T, combine func(T, T) T) T {
	result := slice[0]
	for _, v := range slice[1:] {
		result = combine(result, v)
	}
	return result
}
