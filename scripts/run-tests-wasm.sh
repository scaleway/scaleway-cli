SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../wasm" && pwd)"

TEST_FLAGS=(
  --run
)

if [[ "${TESTER}" == "true" ]]; then
  echo "building"
  GOOS=js GOARCH=wasm go build -o "$ROOT_DIR/cliTester.wasm" ./cmd/scw-wasm-tester
  TEST_FLAGS+=(-t "With test environment")
else
  TEST_FLAGS+=(-t "With wasm CLI")
fi

cd $ROOT_DIR
pnpm install
pnpm test -- "${TEST_FLAGS[@]}"
