SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../wasm" && pwd)"

cd $ROOT_DIR
pnpm install
pnpm test -- --run
