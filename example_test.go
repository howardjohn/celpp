package celpp_test

import (
	"fmt"
	"github.com/howardjohn/celpp"
	"github.com/howardjohn/celpp/macros"
)

func Example_builtinIndex() {
	parser, _ := celpp.New(macros.All...)
	output, _ := parser.Process("self.index(x, z, b)")
	fmt.Println(output)
	// Output: (has(self.x) && has(self.x.z) && has(self.x.z.b)) ? self.x.z.b : null
}

func Example_builtinDefault() {
	parser, _ := celpp.New(macros.All...)
	output, _ := parser.Process("default(self.x, 'DEF')")
	fmt.Println(output)
	// Output: has(self.x) ? self.x : "DEF"
}

func Example_builtinOneof() {
	parser, _ := celpp.New(macros.All...)
	output, _ := parser.Process("oneof(self.x, self.y, self.z)")
	fmt.Println(output)
	// Output: (has(self.x) ? 1 : 0) + (has(self.y) ? 1 : 0) + (has(self.z) ? 1 : 0) <= 1
}