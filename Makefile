VLANG_FILES := $(wildcard *.v)
WASM_FILES := $(patsubst %.v,%.wasm,$(VLANG_FILES))

# Keep test at the top so that it is default when `make` is called.
# This is used by Travis CI.
coverage.txt:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
view-cover: clean coverage.txt
	go tool cover -html=coverage.txt
test: build
	go test ./...
build: build-wasm
	go build ./...
install: build
	go install ./...
	rm wasm-runtime-example
run:
	wasm-runtime-example
update:
	go get -u ./...
pre-commit: update coverage.txt
	go mod tidy

# Vlang / WASM
build-wasm: $(WASM_FILES)

%.wasm: %.v
	mkdir -p wasmout
	echo "Compiling $< to $@"
	v -b wasm -os wasi $< -o $@