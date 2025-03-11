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

	// Export the contained gojaHandleRequest function
	var gojaHandleRequest func(GojaEnv, *http.Request)
	err = vm.ExportTo(vm.Get("handleRequest"), &gojaHandleRequest)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", gojaMiddleware(env, gojaHandleRequest, handler))
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err)
	}
}

func gojaMiddleware(
	env GojaEnv,
	gojaFn func(GojaEnv, *http.Request),
	del http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call gojaFn
		gojaFn(env, r)

		// Carry on with handler
		del(w, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Always return a duck
	w.Write([]byte("a duck"))
}

// Defines go methods that we pass to Goja, so that Goja can use them
type GojaEnv struct{}

func (GojaEnv) Println(a ...any) {
	fmt.Println(a...)
}

func (GojaEnv) Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}
