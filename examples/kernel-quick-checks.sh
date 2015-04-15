#!/usr/bin/env bash

# I'm a script used to check the stability of new kernels.
#
# I spawn 16 servers, boot them, and perform several checks such as a
# local iperf, some IOs (parallel download of 5GO on a NBD mounted
# disk).
#
# I generate a Markdown report with results

# parameters
if [ $# -ne 1 ]; then
    echo "usage: $0 build-id"
    exit 1
fi

# globals
BUILD_ID=$1
WORKDIR=$(mktemp -d -t minibench)
NB_INSTANCES=16
INSTANCE_NAME='minibench-kernel'

# we expect the following in the environment
#
# - JENKINS_MIRROR
# - SANDBOX_DTB_UUID
if [ -f env.bash ]
then
    source env.bash
fi

# fetch kernel and publish to S3
function prepare {
    echo >&2 '[+] preparing kernel...'
    local kernel=$(wget "${JENKINS_MIRROR}/" -O /dev/stdout 2>/dev/null | sed -e 's/<a href="\([^"]*\)\/">.*/\1/' -e 'tx' -e 'd' -e ':x' | grep "^.*\-${BUILD_ID}$\|^.*\-${BUILD_ID}\-.*$")
    if [ -z "$kernel" ]; then
	echo "can't find $BUILD_ID on jenkins"
	exit 1
    fi

    wget "${JENKINS_MIRROR}/${kernel}/uImage" -O /tmp/uImage 2>/dev/null
    wget "${JENKINS_MIRROR}/${kernel}/dtbs/pimouss-computing.dtb" -O /tmp/dtb 2>/dev/null

    s3cmd put --acl-public /tmp/uImage s3://mxs/uImage-$kernel &>/dev/null
    s3cmd put --acl-public /tmp/dtb s3://mxs/dtb-$kernel &>/dev/null

    s3cmd put --acl-public /tmp/uImage s3://mxs/uImage-sandbox &>/dev/null
    s3cmd put --acl-public /tmp/dtb s3://mxs/dtb-sandbox &>/dev/null

    rm -f /tmp/uImage /tmp/dtb
}

# destroy all existing servers matching name
function cleanup {
    echo >&2 '[+] cleaning up existing servers...'
    for uuid in $(scw ps -a --no-trunc | tail -n +2 | awk '// { print $1, $NF; }' | grep "^.* ${INSTANCE_NAME}\-" | awk '// { print $1; }'); do
	scw stop -t $uuid
    done

    touch $WORKDIR/uuids.txt
    touch $WORKDIR/ips.txt
}

# create $NB_INSTANCES servers with a bootscript pointing to the prepared kernel
function boot {
    echo >&2 "[+] creating $NB_INSTANCES servers..."
    for i in $(eval echo {1..$NB_INSTANCES}); do
	scw create --bootscript $SANDBOX_DTB_UUID --volume 50G --name "$INSTANCE_NAME-$i" Ubuntu_Trusty_14_04_LTS >> $WORKDIR/uuids.txt
    done
    cat $WORKDIR/uuids.txt

    echo >&2 "[+] booting $NB_INSTANCES servers..."
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw start -s $uuid &
	sleep .5
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

    # uname -a
    echo >&2 "[+] report uname"
    echo "## uname"
    echo ""
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'uname -a' | sed 's/\(.*\)/    \1/'
    done
    echo ""

    # iperf
    echo >&2 "[+] report iperf"
    echo "## iperf"
    echo ""
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'iperf -s & sleep 5 ; iperf -c localhost' | sed 's/\(.*\)/    \1/'
    done
    echo ""

    # quick stability check
    echo >&2 "[+] report stability 1st pass"
    echo "## stability"
    echo ""
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'find /usr -type f | xargs md5sum &> /tmp/a'
	scw exec $uuid 'find /usr -type f | xargs cat &> /tmp/megafile'
	scw exec $uuid 'for i in {1..5}; do wget --no-verbose --page-requisites http://ping.online.net/1000Mo.dat -O $i 2>/dev/null & done; wait $(jobs -p)' | sed 's/\(.*\)/    \1/'
    done

    echo >&2 "[+] report stability 2nd pass"
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'find /usr -type f | xargs md5sum &> /tmp/b'
    done

    echo >&2 "[+] report stability 3rd pass"
    for uuid in $(cat $WORKDIR/uuids.txt); do
	scw exec $uuid 'diff /tmp/a /tmp/b' | sed 's/\(.*\)/    \1/'
    done

    echo >&2 "[+] report stability fping"
    echo ""
    fping $(cat $WORKDIR/ips.txt) | sed 's/\(.*\)/    \1/'
    echo ""
}

function main {
    prepare
    cleanup
    boot
    report > report-${BUILD_ID}.md
    echo >&2 "[+] report is at report-${BUILD_ID}.md"
    cleanup
}

main
rm -rf $WORKDIR
