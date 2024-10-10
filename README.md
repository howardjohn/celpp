# CEL Pre-processor

[![Go Reference](https://pkg.go.dev/badge/github.com/howardjohn/celpp.svg)](https://pkg.go.dev/github.com/howardjohn/celpp)

This repo contains a simple CEL pre-processor.
This allows passing in higher-level CEL expressions, and emitting standard ones.

This is done using CEL's own parser and macro system.
The pre-processor allows defining custom macros, expands them, and emits the expanded expression.
Builtin macros are *not* expanded, so the result is not overly bloated.

## Examples

Below shows some examples.
This is an incomplete list, see the [Go pkg documentation](https://pkg.go.dev/github.com/howardjohn/celpp) for an exhaustive list.

Note this is also just the builtin ones; library users can define their own as well.

### Default

Input: `has(self.x) && default(self.y, 0)`

Output: `has(self.x) && (has(self.y) ? self.y : 0)`

### Index
Input: `self.index({}, x, z, b)`

Output: `(has(self.x) && has(self.x.z) && has(self.x.z.b)) ? self.x.z.b : {}`

## Future additions

* Custom operator overloads
* Schema awareness. For instance, rather than `default(self.y, 0)`, can we do `default(self.y)` if we know `y` is an int?
