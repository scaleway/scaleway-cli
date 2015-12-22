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
IMAGE_BOOTSCRIPT=${IMAGE_BOOTSCRIPT:stable}
VOLUME_SIZE=${VOLUME_SIZE:-50GB}
KEY=$(cat ~/.ssh/id_rsa.pub | awk '{ print $1" "$2 }' | tr ' ' '_')

echo "[+] URL of the tarball: ${URL}"
echo "[+] Target name: ${NAME}"


echo "[+] Creating new server in rescue mode with a secondary volume..."
SERVER=$(scw run -d --bootscript=rescue --volume="${VOLUME_SIZE}" --env="AUTHORIZED_KEY=${KEY}" --name="image-writer-${NAME}" 1GB)
echo "[+] Server created: ${SERVER}"


echo "[+] Booting..."
scw exec -w "${SERVER}" 'uname -a'
echo "[+] Server is booted"


if [ -n "${SCW_GATEWAY_HTTP_PROXY}" ]; then
    echo "[+] Configuring HTTP proxy"
    # scw exec "${SERVER}" "echo proxy=${SCW_GATEWAY_HTTP_PROY} >> .curlrc"
    (
        set +x
        scw exec "${SERVER}" "echo http_proxy = ${SCW_GATEWAY_HTTP_PROXY} >> .wgetrc" >/dev/null 2>/dev/null || (echo "Failed to configure HTTP proxy"; exit 1) || exit 1
    )
fi


echo "[+] Formating and mounting /dev/nbd1..."
scw exec "${SERVER}" 'mkfs.ext4 /dev/nbd1 && mount /dev/nbd1 /mnt'
echo "[+] /dev/nbd1 formatted in ext4 and mounted on /mnt"


echo "[+] Download tarball and write it to /dev/nbd1"
scw exec "${SERVER}" "wget -qO - ${URL} | tar -C /mnt/ -xf - && sync"
echo "[+] Tarball extracted on /dev/nbd1"


echo "[+] Stopping the server"
scw stop "${SERVER}" >/dev/null
scw wait "${SERVER}"
echo "[+] Server stopped"


echo "[+] Creating a snapshot of nbd1"
SNAPSHOT=$(scw commit --volume=1 "${SERVER}" "${SNAPSHOT_NAME}")
echo "[+] Snapshot ${SNAPSHOT} created"


echo "[+] Creating an image based of the snapshot"
if [ -n "${IMAGE_ARCH}" ]; then
    IMAGE=$(scw tag --arch="${IMAGE_ARCH}" --bootscript="${IMAGE_BOOTSCRIPT}" "${SNAPSHOT}" "${IMAGE_NAME}")
else
    IMAGE=$(scw tag --bootscript="${IMAGE_BOOTSCRIPT}" "${SNAPSHOT}" "${IMAGE_NAME}")
fi
echo "[+] Image created: ${IMAGE}"


echo "[+] Deleting temporary server"
scw rm "${SERVER}" >/dev/null
echo "[+] Server deleted"
