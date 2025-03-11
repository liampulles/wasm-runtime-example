package main

import (
	_ "embed"
	"fmt"

	"github.com/dop251/goja"
)

//go:embed add.js
var addJS []byte

// Compile the javascript
var addProg = goja.MustCompile("add", string(addJS), true)

func main() {
	// Create goja runtime
	vm := goja.New()

	// Load the add program into the vm (we run the script, which loads declared funcs)
	_, err := vm.RunProgram(addProg)
	if err != nil {
		panic(err)
	}

	// Export the contained add function
	var add func(int, int) int
	err = vm.ExportTo(vm.Get("add"), &add)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d + %d = %d\n", 1, 2, add(1, 2))
}
