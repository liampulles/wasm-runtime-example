package main

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/dop251/goja"
)

//go:embed middleware.js
var addJS []byte

// Compile the javascript
var addProg = goja.MustCompile("middleware", string(addJS), true)

func main() {
	// Create goja runtime
	vm := goja.New()

	// Load the add program into the vm (we run the script, which loads declared funcs)
	_, err := vm.RunProgram(addProg)
	if err != nil {
		panic(err)
	}

	// Instantiate our go env which we will pass through
	env := GojaEnv{}

	// Export the contained handleRequest function
	var handleRequest func(GojaEnv, http.Request)
	err = vm.ExportTo(vm.Get("handleRequest"), &handleRequest)
	if err != nil {
		panic(err)
	}

}

// Defines go methods that we pass to Goja, so that Goja can use them
type GojaEnv struct{}

func (GojaEnv) Println(a ...any) {
	fmt.Println(a...)
}

func (GojaEnv) Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}
