package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed middleware/log.wasm
var logWasm []byte
var logMiddleware RequestTapMiddleware

// Request tapping middleware can just read the request, but not modify it or
// stop it. Suitable for logging and monitoring.
type RequestTapMiddleware func(http.Request)

// This is close to what the web assembly function must accept.
type wasmRequestTapMiddleware func(path string, method string, body string)

func mustLoadRequestTapMiddleware(ctx context.Context, r wazero.Runtime, wasm []byte) RequestTapMiddleware {
	// Read the wasm bytes
	mod, err := r.InstantiateWithConfig(ctx, wasm, wazero.NewModuleConfig())
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	// Get the wasm function
	wasmFn := mod.ExportedFunction("middleware")
	if wasmFn == nil {
		err := errors.New("middleware() does not seem to be exported from the wasm - is it exported? Suggest using wasm-tools to check")
		log.Panic(err)
	}

	// Wrap the bridge func. This will need to allocate and free memory.
	var bridgeFn wasmRequestTapMiddleware = func(path, method, body string) {
		wasmFn.Call(ctx, wazero, method, body)
	}
}

func wasmStr(in string) uint64

func main() {
	ctx := context.Background()

	// Setup the WASI runtime
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Load middleware

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":9090", nil)

	// Call the `add` function and print the results to the console.
	add := mod.ExportedFunction("add")
	results, err := add.Call(ctx, 1, 2)
	if err != nil {
		log.Panicf("failed to call add: %v", err)
	}

	fmt.Printf("%d + %d = %d\n", 1, 2, results[0])
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
