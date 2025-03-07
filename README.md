# wasm-runtime-example
Example app to see how to execute WASM from Go, for use as HTTP middleware

The middleware code is written in Vlang (the compiled to WASI [Preview 1] compatible WASM).

## Helpful

Install wasm-tools: https://github.com/bytecodealliance/wasm-tools
It is very useful for inspecting compiled WASM.

### Check what functions a .wasm file contains

```shell
wasm-tools dump wasmout/add.wasm | grep export
```

## Get a .wat file from a .wasm

```shell
wasm-tools print wasmout/add.wasm -o wasmout/add.wat
```
