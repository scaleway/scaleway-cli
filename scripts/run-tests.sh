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

  color green "  -u, --update"
  echo -e "\tUpdate goldens during integration tests."

  color green "  -g, --update-goldens"
  echo -e "\tUpdate goldens during integration tests."

  color green "  -c, --update-cassettes"
  echo -e "\tRecord cassettes during integration tests. This requires a valid Scaleway token in your environment."

  color green "  -D, --debug"
  echo -e "\tEnable CLI debug mode."

  color green "  -h, --help"
  echo -e "\tDisplay this help."

  echo ""
  exit $1;
}

SCW_DEBUG="false"
OPT_UPDATE_GOLDENS="false"
OPT_UPDATE_CASSETTES="false"

##
# Parse arguments
##
while [[ $# > 0 ]]
do
  case "$1" in
    -u|--update)
      OPT_UPDATE_GOLDENS="true"
      OPT_UPDATE_CASSETTES="true"
      ;;
    -g|--update-goldens)
      OPT_UPDATE_GOLDENS="true"
      ;;
    -c|--update-cassettes)
      OPT_UPDATE_CASSETTES="true"
      ;;
    -D|--debug)
      SCW_DEBUG="true"
      ;;
	-h|--help) usage
  esac
  shift
done

SCW_DEBUG=$SCW_DEBUG CLI_UPDATE_GOLDENS=$OPT_UPDATE_GOLDENS CLI_UPDATE_CASSETTES=$OPT_RECORD_CASSETTES go test -v $ROOT_DIR/...
