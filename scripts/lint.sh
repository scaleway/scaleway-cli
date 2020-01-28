#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

##
# Colorize output
##
function color() {
  case $1 in
    yellow) echo -e -n "\033[33m"   ;;
    green)  echo -e -n "\033[32m"   ;;
    red)    echo -e -n "\033[0;31m" ;;
  esac
  echo "$2"
  echo -e -n "\033[0m"
}

##
# Print Usage
##
function usage() {
  color yellow "Usage:"
  echo "  $SCRIPT_DIR [OPTIONS]"
  echo ""

  color yellow "Options:"

  color green "  -w, --write"
  echo -e "\tFix found issues (if it's supported by the linter)."

  color green "  --list"
  echo -e "\tList current linters configuration."

  color green "  -h, --help"
  echo -e "\tDisplay this help."

  echo ""
  exit "$1";
}

OPT_CMD="run"
OPT_FLAGS=""
OPT_DIRS="$ROOT_DIR/..."

##
# Parse arguments
##
while [[ $# -gt 0 ]]
do
  case "$1" in
	  -h|--help)
	    usage 0 ;;

    -w|--write)
      OPT_CMD="run"
      OPT_FLAGS+=" --fix" ;;

    --list)
      OPT_DIRS=""
      OPT_CMD="linters" ;;

    -v|--verbose)
      OPT_FLAGS+=" -v" ;;

    *)
      color red "Unkown argument '$1'"
      echo
      usage 1 ;;
  esac
  shift
done

##
# Check golangci-lint command existence
##
if [ ! -x "$(command -v golangci-lint)" ];
then
  echo "golangci-lint is not installed"
  echo "On macOS, you can run: brew install golangci/tap/golangci-lint"
  echo "On other systems, refer to installation instructions: https://github.com/golangci/golangci-lint#install"
  exit 1
fi

##
# Execute golangci-lint command
##
golangci-lint $OPT_CMD $OPT_FLAGS $OPT_DIRS
