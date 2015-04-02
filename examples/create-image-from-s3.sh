#!/bin/bash

set -e
URL="${1}"

if [ -z "${1}" ]; then
    echo "usage: $(basename ${0}) <url>"
    echo ""
    echo "examples:"
    echo "  - $(basename ${0}) http://test-images.fr-1.storage.online.net/ocs-distrib-ubuntu-trusty.tar"
    echo "  - VOLUME=20GB $(basename ${0}) http://test-images.fr-1.storage.online.net/ocs-distrib-ubuntu-trusty.tar"
    exit 1
fi

# FIXME: add usage

set -e

NAME=$(basename "${URL}")
NAME=${NAME%.*}-$(date +%Y-%m-%d_%H:%M)
VOLUME_SIZE=${VOLUME_SIZE:-50GB}


echo "[+] URL of the tarball: ${URL}"
echo "[+] Target name: ${NAME}"


echo "[+] Creating new server in rescue mode with a secondary volume..."
SERVER=$(onlinelabs create 1GB --bootscript=rescue --volume="${VOLUME_SIZE}" --name="image-writer-${NAME}")
echo "[+] Server created: ${SERVER}"


echo "[+] Booting..."
onlinelabs start --sync "${SERVER}" >/dev/null
IP=$(onlinelabs inspect "${SERVER}" -f .server.public_ip.address)
onlinelabs exec --insecure "${SERVER}" 'uname -a'
echo "[+] SSH is ready (${IP})"


echo "[+] Formating and mounting /dev/nbd1..."
onlinelabs exec "${SERVER}" 'service xnbd-common stop && service xnbd-common start && mkfs.ext4 /dev/nbd1 && mount /dev/nbd1 /mnt'
echo "[+] /dev/nbd1 formatted in ext4 and mounted on /mnt"


echo "[+] Download tarball from S3 and write it to /dev/nbd1"
onlinelabs exec "${SERVER}" "wget -qO - ${URL} | tar -C /mnt/ -xf - && sync"
echo "[+] Tarball extracted on /dev/nbd1"


echo "[+] Stopping the server"
onlinelabs stop "${SERVER}"
onlinelabs wait "${SERVER}"
echo "[+] Server stopped"


echo "[+] Creating a snapshot of nbd1"
SNAPSHOT=$(onlinelabs commit "${SERVER}" --volume=1 --name="${NAME}")
echo "[+] Snapshot ${SNAPSHOT} created"


echo "[+] Creating an image based of the snapshot"
IMAGE=$(onlinelabs tag "${SNAPSHOT}" "${NAME}")
echo "[+] Image created: ${IMAGE}"


echo "[+] Deleting temporary server"
onlinelabs rm "${SERVER}"
echo "[+] Server deleted"
