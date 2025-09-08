#/bin/env bash

MAX_BINARY_SIZE=53000000

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function usage() {
  echo "Usage:"
  echo "  $SCRIPT_DIR [BINARY_FILE]"

  echo ""
  exit "$1";
}

if [[ $# -ne 1 ]]; then
  usage
fi

BINARY_SIZE=$(stat -c %s $1)

echo "Binary size: $BINARY_SIZE"

if (( $BINARY_SIZE > $MAX_BINARY_SIZE )); then
  echo "Size is greater than limit"
  exit 1
else
  echo "Size is valid"
fi
