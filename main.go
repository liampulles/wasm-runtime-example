package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasmout/add.wasm
var addWasm []byte

func main() {
	ctx := context.Background()

	// Create a new WebAssembly Runtime, close it later
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI
	wasi_snapshot_preview1.MustInstantiate(ctx, r)
	// Load add.wasm
	addMod, err := r.InstantiateWithConfig(ctx, addWasm, wazero.NewModuleConfig())
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	// Call the `add` function and print the results to the console.
	add := addMod.ExportedFunction("add")
	results, err := add.Call(ctx, 1, 2)
	if err != nil {
		log.Panicf("failed to call add: %v", err)
	}

	fmt.Printf("%d + %d = %d\n", 1, 2, results[0])
}
