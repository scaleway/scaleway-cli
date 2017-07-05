#!/bin/bash

set -e
URL="${1}"

if [ -z "${1}" ]; then
    echo "usage: $(basename ${0}) <url>"
    echo ""
    echo "examples:"
    echo "  - $(basename ${0}) http://test-images.fr-1.storage.online.net/scw-distrib-ubuntu-trusty.tar"
    echo "  - VOLUME_SIZE=50GB $(basename ${0}) http://test-images.fr-1.storage.online.net/scw-distrib-ubuntu-trusty.tar"
    exit 1
fi


NAME=$(basename "${URL}")
SNAPSHOT_NAME=${NAME%.*}-$(date +%Y-%m-%d_%H:%M)
IMAGE_NAME=${IMAGE_NAME:-$SNAPSHOT_NAME}
SCW_COMMERCIAL_TYPE=${SCW_COMMERCIAL_TYPE:-VC1S}
VOLUME_SIZE=${VOLUME_SIZE:-50GB}
SCW_TARGET_ARCH=x86_64
if [ "$SCW_COMMERCIAL_TYPE" = "C1" ]; then
    SCW_TARGET_ARCH=arm
fi
KEY=$(cat ~/.ssh/id_rsa.pub | awk '{ print $1" "$2 }' | tr ' ' '_')

echo "[+] URL of the tarball: ${URL}"
echo "[+] Target name: ${NAME}"


echo "[+] Creating new server in rescue mode with a secondary volume..."
SERVER=$(SCW_TARGET_ARCH="$SCW_TARGET_ARCH" scw run --commercial-type="${SCW_COMMERCIAL_TYPE}" -d --env="AUTHORIZED_KEY=${KEY} boot=none INITRD_DROPBEAR=1" --name="image-writer-${NAME}" "${VOLUME_SIZE}")
echo "[+] Server created: ${SERVER}"


echo "[+] Booting..."
scw exec -w "${SERVER}" 'uname -a'
echo "[+] Server is booted"


if [ -n "${SCW_GATEWAY_HTTP_PROXY}" ]; then
    echo "[+] Configuring HTTP proxy"
    # scw exec "${SERVER}" "echo proxy=${SCW_GATEWAY_HTTP_PROY} >> .curlrc"
    (
        set +x
        scw exec "${SERVER}" "echo export http_proxy=${SCW_GATEWAY_HTTP_PROXY} > /proxy-env" >/dev/null 2>/dev/null || (echo "Failed to configure HTTP proxy"; exit 1) || exit 1
    )
fi


echo "[+] Formating and mounting disk..."
# FIXME: make disk dynamic between /dev/vda and /dev/nbd1
scw exec "${SERVER}" '/sbin/mkfs.ext4 /dev/vda && mkdir -p /mnt && mount /dev/vda /mnt'
echo "[+] /dev/nbd1 formatted in ext4 and mounted on /mnt"


echo "[+] Download tarball and write it to /mnt"
scw exec "${SERVER}" "touch /proxy-env; . /proxy-env; wget -qO - ${URL} | tar -C /mnt/ -xf - && sync"
echo "[+] Tarball extracted on disk"


echo "[+] Stopping the server"
scw stop "${SERVER}" >/dev/null
scw wait "${SERVER}"
echo "[+] Server stopped"


echo "[+] Creating a snapshot of disk 1"
SNAPSHOT=$(scw commit --volume=0 "${SERVER}" "${SNAPSHOT_NAME}")
echo "[+] Snapshot ${SNAPSHOT} created"


echo "[+] Creating an image based of the snapshot"
tag_args=()

[ -n "${IMAGE_ARCH}" ]       && tag_args+=(--arch="${IMAGE_ARCH}")
[ -n "${IMAGE_BOOTSCRIPT}" ] && tag_args+=(--bootscript="${IMAGE_BOOTSCRIPT}")

IMAGE=$(scw tag "${tag_args[@]}" "${SNAPSHOT}" "${IMAGE_NAME}")

echo "[+] Image created: ${IMAGE}"


echo "[+] Deleting temporary server"
scw rm "${SERVER}" >/dev/null
echo "[+] Server deleted"
