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
	UnrollMap,
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
	// Usage: `self.index({}, x, z, b)`
	// This does a nil-safe traversal of self.x.z.b. If any field is nil, the first arg ({} here) is returned.
	Index = cel.ReceiverVarArgMacro("index",
		func(mef cel.MacroExprFactory, base celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
			if len(args) < 2 {
				return nil, mef.NewError(base.ID(), "index requires at least 2 arg")
			}
			zero := args[0]
			checks := []celast.Expr{}
			for i := range args[1:] {
				next := mef.NewCall(operators.Has, selects(mef, base, args[1:i+2]...))
				checks = append(checks, next)
			}
			check := foldl(checks, func(l, r celast.Expr) celast.Expr {
				return mef.NewCall(operators.LogicalAnd, l, r)
			})

			final := mef.NewCall(operators.Conditional, check, selects(mef, base, args[1:]...), zero)
			return final, nil
		})

	// UnrollMap provides a specialized macro to unroll a loop. This works around bugs in older Kubernetes versions cost estimation.
	// Usage: `self.unrollmap(0, 3, x, x.matches.size())`
	// Instead of `self.map(x, x.matches.size())`.
	// The unrolling will *always* create a list of N elements (3 in above example), unlike map.
	// Since the input list may not be that size, you also need a zero-value to fill in gaps (0 in above example)
	UnrollMap = cel.ReceiverVarArgMacro("unrollmap",
		func(mef cel.MacroExprFactory, base celast.Expr, args []celast.Expr) (celast.Expr, *cel.Error) {
			if len(args) == 0 {
				return nil, mef.NewError(base.ID(), "unrollmap requires at 4 args")
			}
			zero := args[0]
			count := (args[1].AsLiteral().Value()).(int64)
			varname := args[2]
			expr := args[3]
			items := []celast.Expr{}
			for n := range count {
				sizeCheck := mef.NewCall(operators.Greater,
					mef.NewCall("size", base),
					mef.NewLiteral(types.Int(n)),
				)
				expr := mef.NewCall(operators.Index,
					mef.NewMemberCall(operators.Map,
						mef.NewList(mef.NewCall(operators.Index, base, mef.NewLiteral(types.Int(n)))),
						varname,
						expr,
					),
					mef.NewLiteral(types.Int(0)),
				)
				clause := mef.NewCall(operators.Conditional,
					sizeCheck,
					expr,
					zero,
				)
				items = append(items, clause)
			}
			return mef.NewList(items...), nil
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
