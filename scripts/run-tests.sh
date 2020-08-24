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

  color green "  -r, --run <regex>"
  echo -e "\tRun a specific test or set of tests matching the given regex. Similar to the '-run' Go test flag."

  color green "  -u, --update"
  echo -e "\tUpdate goldens and record cassettes during integration tests."

  color green "  -g, --update-goldens"
  echo -e "\tUpdate goldens during integration tests."

  color green "  -c, --update-cassettes"
  echo -e "\tRecord cassettes during integration tests. Warning: a valid Scaleway token is required in your environment in order to record cassettes."

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
OPT_RUN_SCOPE=""

##
# Parse arguments
##
while [[ $# -gt 0 ]]
do
  case "$1" in
    -r|-run|--run) # keeping -run as this is the standard Go flag for this
      shift
      OPT_RUN_SCOPE="$1"
      ;;
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

# Test cache is cleansed before updating cassettes in order to force recreate them.
if [[ ${OPT_UPDATE_CASSETTES} ]] ; then
  go clean -testcache
fi

# Remove golden if they are being updated, and all tests are being run
if [[ ${OPT_UPDATE_GOLDENS} == "true" ]] && [[ -z ${OPT_RUN_SCOPE} ]]; then
  # We ignore OS specific goldens
  find . -type f ! -name '*windows*.golden' ! -name '*darwin*.golden' ! -name '*linux*.golden' -name "*.golden" -exec rm -f {} \;
fi

SCW_DEBUG=$SCW_DEBUG CLI_UPDATE_GOLDENS=$OPT_UPDATE_GOLDENS CLI_UPDATE_CASSETTES=$OPT_UPDATE_CASSETTES go test -v $ROOT_DIR/... -timeout 20m -run=$OPT_RUN_SCOPE
