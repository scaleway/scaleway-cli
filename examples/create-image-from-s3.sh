#!/bin/bash

set -e
URL="${1}"

if [ -z "${1}" ]; then
    echo "usage: $(basename ${0}) <url>"
    echo ""
    echo "examples:"
    echo "  - $(basename ${0}) http://test-images.fr-1.storage.online.net/ocs-distrib-ubuntu-trusty.tar"
    echo "  - VOLUME=20000000000 $(basename ${0}) http://test-images.fr-1.storage.online.net/ocs-distrib-ubuntu-trusty.tar"
    exit 1
fi

# FIXME: add usage

NAME=$(basename "${URL}")
NAME=${NAME%.*}-$(date +%Y-%m-%d_%H:%M)
VOLUME_SIZE=${VOLUME_SIZE:-50000000000}  # 50GB


echo "[+] URL of the tarball: ${URL}"
echo "[+] Target name: ${NAME}"


echo "[+] Creating new server in rescue mode with a secondary volume..."
SERVER=$(onlinelabs create trusty --bootscript=rescue --volume=${VOLUME_SIZE} --name="image-writer-${NAME}")
echo "[+] Server created: ${SERVER}"


echo "[+] Booting..."
onlinelabs start "${SERVER}" >/dev/null


# FIXME: wait for state to be "running"


echo "[+] Waiting for SSH to be available"
until onlinelabs exec "${SERVER}" --insecure -- exit 0; do sleep 5; done &>/dev/null
IP=$(onlinelabs inspect "${SERVER}" -f .server.public_ip.address)
echo "[+] SSH is ready (${IP})"


echo "[+] Formating /dev/nbd1..."
onlinelabs exec "${SERVER}" -- 'service xnbd-common stop && service xnbd-common start && mkfs.ext4 /dev/nbd1'
echo "[+] /dev/nbd1 formatted in ext4"


echo "[+] Mounting /dev/nbd1"
onlinelabs exec "${SERVER}" -- mount /dev/nbd1 /mnt
echo "[+] /dev/nbd1 mounted on /mnt"


echo "[+] Download tarball from S3 and write it to /dev/nbd1"
onlinelabs exec "${SERVER}" -- "wget -qO - ${URL} | tar -C /mnt/ -xf - && sync"
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


# FIXME: cleaning
