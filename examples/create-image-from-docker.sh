#!/bin/bash

set -e
REPO="${1}"

if [ -z "${1}" ]; then
    echo "usage: $(basename ${0}) <repo>"
    echo ""
    echo "examples:"
    echo "  - $(basename ${0}) armbuild/scw-app-docker:latest"
    echo "  - VOLUME_SIZE=50GB $(basename ${0}) armbuild/scw-app-docker"
    exit 1
fi

set -e

NAME=$(basename "${REPO}")
NAME=${NAME%.*}-$(date +%Y-%m-%d_%H:%M)
VOLUME_SIZE=${VOLUME_SIZE:-50GB}


echo "[+] REPO: ${REPO}"
echo "[+] Target name: ${NAME}"


echo "[+] Creating new server in rescue mode with a secondary volume..."
SERVER=$(scw create --volume="${VOLUME_SIZE}" --name="image-writer-${NAME}" docker)
echo "[+] Server created: ${SERVER}"


echo "[+] Booting..."
scw start --wait --timeout=600 "${SERVER}" >/dev/null
#IP=$(scw inspect -f .server.public_ip.address "${SERVER}")
#echo "[+] SSH is ready (${IP})"
echo "[+] Server is booted"
scw exec "${SERVER}" 'uname -a'


echo "[+] Formating and mounting /dev/nbd1..."
scw exec "${SERVER}" 'mkfs.ext4 /dev/nbd1 && mount /dev/nbd1 /mnt'
echo "[+] /dev/nbd1 formatted in ext4 and mounted on /mnt"


echo "[+] Download image"
scw exec "${SERVER}" docker pull "${REPO}"
echo "[+] Image downloaded"

echo "[+] Export image to /dev/nbd1"
scw exec "${SERVER}" docker run --name tmp --entrypoint /dontexists "${REPO}" 2>/dev/null || true
scw exec "${SERVER}" 'docker export tmp > image.tar'
scw exec "${SERVER}" tar -C /mnt -xf image.tar
#scw exec "${SERVER}" rm -f image.tar
scw exec "${SERVER}" rm -f /mnt/.dockerenv /mnt/.dockerinit
scw exec "${SERVER}" chmod 1777 /mnt/tmp
scw exec "${SERVER}" chmod 755 /mnt/etc /mnt/usr /mnt/usr/local /mnt/usr/sbin
scw exec "${SERVER}" chmod 555 /mnt/sys
scw exec "${SERVER}" chmod 700 /mnt/root
scw exec "${SERVER}" mv /mnt/etc/hosts.default /mnt/etc/hosts || true
scw exec "${SERVER}" sync
echo "[+] Image exported"


echo "[+] Stopping the server"
scw stop "${SERVER}"
scw wait "${SERVER}"
echo "[+] Server stopped"


echo "[+] Creating a snapshot of nbd1"
SNAPSHOT=$(scw commit --volume=1 "${SERVER}" "${NAME}")
echo "[+] Snapshot ${SNAPSHOT} created"


echo "[+] Creating an image based of the snapshot"
IMAGE=$(scw tag "${SNAPSHOT}" "${NAME}")
echo "[+] Image created: ${IMAGE}"


echo "[+] Deleting temporary server"
scw rm "${SERVER}"
echo "[+] Server deleted"
