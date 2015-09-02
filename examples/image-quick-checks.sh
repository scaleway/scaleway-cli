#!/usr/bin/env bash

# I'm a script used to check some basic operations with images.

# Params
if [ $# -ne 1 ]; then
    echo "usage: $0 image-id"
    exit 1
fi

# Globals

IMAGE_UUID=$1
IMAGE_NAME=""
SERVER_UUID=""
SERVER_NAME=""
WORKDIR=$(mktemp -d 2>/dev/null || mktemp -d -t /tmp)

# Printing helpers

function _log_msg {
    echo -ne "${*}" >&2
}

function einfo {
    _log_msg "\033[1;36m>>> \033[0m${*}\n"
}

function echeck {
    einfo $@
    echo ""
}

function esuccess {
    _log_msg "\033[1;32m>>> \033[0m${*}\n"
}

function ewarn {
    _log_msg "\033[1;33m>>> \033[0m${*}\n"
}

function eerror {
    _log_msg "\033[1;31m>>> ${*}\033[0m\n"
}

function eedie_on_error {
    failed=$1
    if [ $failed -ne 0 ]
    then
	eerror $2
	cleanup
	exit 1
    fi
}

# Super helper to print server output to stdout

function watch_server {
    einfo "Attaching to server $SERVER_UUID"
    mkfifo ${WORKDIR}/stop_watching_server
    scw attach $SERVER_UUID 2>/dev/null &
    killme=$!
    sleep 5
    (
	cat ${WORKDIR}/stop_watching_server
	kill -9 $killme
	echo ""
    ) &>/dev/null &
}

function stop_watching {
    echo "" > ${WORKDIR}/stop_watching_server
    rm -f ${WORKDIR}/stop_watching_server
    sleep 1
}

# Checks

function check_image {
    echeck "Checking image..."
    IMAGE_NAME=$(scw inspect -f '{{ .Name }}' $IMAGE_UUID 2> /dev/null)
    eedie_on_error $? "Unable to find image behind $IMAGE_UUID"
    esuccess "Image name is $IMAGE_NAME"
    SERVER_NAME="qa-image-$$"
}

function check_create {
    echeck "Checking server creation with image..."
    SERVER_UUID=$(scw create --name "$SERVER_NAME" $IMAGE_UUID)
    eedie_on_error $? "Unable to create server with image $IMAGE_NAME"
    esuccess "Created server $SERVER_UUID with image $IMAGE_NAME"
}

function check_boot {
    echeck "Checking server boot (can take up to 300 seconds)..."
    watch_server $SERVER_UUID
    scw start -w -T 300 $SERVER_UUID &> /dev/null
    eedie_on_error $? "Unable to boot server $SERVER_UUID"
    stop_watching
    esuccess "Server $SERVER_UUID properly boots"
}

function check_hostname {
    echeck "Checking server hostname..."
    HOSTNAME=$(scw exec $SERVER_UUID hostname)
    eedie_on_error $? "Unable to fetch hostname from server $SERVER_UUID"
    if [ "$HOSTNAME" -ne "$SERVER_NAME" ]
    then
	ewarn "Server hostname is *NOT* properly set (expected=${SERVR_NAME}, found=${HOSTNAME})"
    else
	esuccess "Server hostname is properly set"
    fi
}

function check_reboot {
    echeck "Checking reboot..."
    watch_server
    scw exec $SERVER_UUID reboot &> /dev/null
    scw exec -T 300 -w $SERVER_UUID uptime
    stop_watching
    eedie_on_error $? "Unable to reboot server $SERVER_NAME"
    esuccess "Server $SERVER_NAME properly reboots"
}

# Cleanup

function cleanup {
    if [ "$SERVER_UUID" != "" ]
    then
	einfo "Cleaning up server $SERVER_UUID"
	scw stop -t $SERVER_UUID &> /dev/null
	scw rm $SERVER_UUID &> /dev/null
	SERVER_UUID=""
    fi
    stop_watching
    einfo "Cleaning up temporary directory $WORKDIR"
    rm -rf $WORKDIR
    kill -9 $(jobs -p)
    einfo "Killed all survivor processes"
}

# Main

function main {
    check_image
    check_create
    check_boot
    check_hostname
    check_reboot
    cleanup
}

main
