SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../wasm" && pwd)"

echo "building WASM binaries"
GOOS=js GOARCH=wasm go build -o "$ROOT_DIR/cli.wasm" ./cmd/scw-wasm
GOOS=js GOARCH=wasm go build -o "$ROOT_DIR/cliTester.wasm" ./cmd/scw-wasm-tester

cd "$ROOT_DIR"
pnpm install
pnpm test -- --run
