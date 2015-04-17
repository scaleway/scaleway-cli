#!/usr/bin/env bash

# I'm a script used to check the state of images.

# parameters
if [ $# -ne 1 ]; then
    echo "usage: $0 image-id"
    exit 1
fi

IMAGE_ID=$1
NB_INSTANCES=16
INSTANCE_NAME='check-image'

# destroy all existing servers matching name
function cleanup {
    echo >&2 '[+] cleaning up existing servers...'
    for uuid in $(scw ps -a --no-trunc | tail -n +2 | awk '// { print $1, $NF; }' | grep "^.* ${INSTANCE_NAME}\-" | awk '// { print $1; }'); do
	scw stop -t $uuid
    done
    
    touch $WORKDIR/uuids.txt
    touch $WORKDIR/ips.txt
}

# create $NB_INSTANCES servers using the image
function boot {
    echo >&2 "[+] creating $NB_INSTANCES servers..."
    for i in $(eval echo {1..$NB_INSTANCES}); do
	scw create --bootscript $SANDBOX_DTB_UUID --volume 50G --volume 1G --name "$INSTANCE_NAME-$i" $IMAGE_ID >> $WORKDIR/uuids.txt
    done
    cat $WORKDIR/uuids.txt

    echo >&2 "[+] booting $NB_INSTANCES servers..."
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw start -s --boot-timeout=120 --ssh-timeout=120 $uuid &
    done
    wait `jobs -p`

    echo >&2 "[+] fetching IPs..."
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw inspect $uuid | grep address | awk '// { print $2; }' | tr -d '"' | awk '// { print $1; }' >> $WORKDIR/ips.txt
    done
}

# run several tests and output a Markdown report
function report {
    # status
    echo >&2 "[+] report status"
    echo "## Status of instances"
    echo ""
    NB_INSTANCES_OK=$(wc -l $WORKDIR/ips.txt | awk '// { print $1; }')
    echo "- $NB_INSTANCES_OK / $NB_INSTANCES have correctly booted"
    echo ""

    # fping
    echo >&2 "[+] report fping"
    echo "## fping"
    echo ""
    fping $(cat $WORKDIR/ips.txt) | sed 's/\(.*\)/    \1/'
    echo ""

    # reboot
    echo >&2 "[+] reboot"
    echo "## reboot"
    echo ""
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'systemctl reboot ; sleep 10 ; reboot' &
    done
    echo ""

    sleep 120
    kill -9 `jobs -p`

    # fping
    echo >&2 "[+] report fping after reboot"
    echo "## fping after reboot"
    echo ""
    fping $(cat $WORKDIR/ips.txt) | sed 's/\(.*\)/    \1/'
    echo ""

    # fping
    echo >&2 "[+] uptime"
    echo "## uptime"
    echo ""
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'uptime'
    done
    echo ""
}

function main {
    cleanup
    boot
    report
}

main
