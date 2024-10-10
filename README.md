# CEL Pre-processor

This repo contains a simple CEL pre-processor.
This allows passing in higher-level CEL expressions, and emitting standard ones.

This is done using CEL's own parser and macro system.
The pre-processor allows defining custom macros, expands them, and emits the expanded expression.
Builtin macros are *not* expanded, so the result is not overly bloated.

## Examples

Input: `has(self.x) && default(self.y, 0)`
Output: `has(self.x) && (has(self.y) ? self.y : 0)`

Input: `self.index(x, z, b)`
Output: `(has(self.x) && has(self.x.z) && has(self.x.z.b)) ? self.x.z.b : null`

## Status

This is a complete WIP at this point, but I will likely pick up work in the future.
