#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

OPT_VERBOSE=""

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

  color green "  --update"
  echo -e "\tUpdate goldens during integration tests."

  color green "  -h, --help"
  echo -e "\tDisplay this help."

  echo ""
  exit $1;
}

OPT_UPDATE_GOLDENS="false"

##
# Parse arguments
##
while [[ $# > 0 ]]
do
  case "$1" in
	-h|--help) usage ;;
    --update)
      OPT_UPDATE_GOLDENS="true"
  esac
  shift
done

UPDATE_GOLDEN=$OPT_UPDATE_GOLDENS go test $ROOT_DIR/...
