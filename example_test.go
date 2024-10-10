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

func Example_builtinUnrollmap() {
	parser, _ := celpp.New(macros.All...)
	output, _ := parser.Process("self.unrollmap(0, 3, x, x.matches.size()).sum()<=128")
	fmt.Println(output)
	// Output: [(size(self) > 0) ? ([self[0]].map(x, x.matches.size())[0]) : 0, (size(self) > 1) ? ([self[1]].map(x, x.matches.size())[0]) : 0, (size(self) > 2) ? ([self[2]].map(x, x.matches.size())[0]) : 0].sum() <= 128
}
