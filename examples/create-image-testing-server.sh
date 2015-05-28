#!/bin/bash

set -e
URL="${1}"

if [ -z "${1}" ]; then
    echo "usage: $(basename ${0}) <url>"
    echo ""
    echo "examples:"
    echo "  - $(basename ${0}) http://test-images.fr-1.storage.online.net/scw-distrib-ubuntu-trusty.tar"
    echo "  - VOLUME=20GB $(basename ${0}) http://test-images.fr-1.storage.online.net/scw-distrib-ubuntu-trusty.tar"
    exit 1
fi

# FIXME: add usage

NAME=$(basename "${URL}")
NAME=${NAME%.*}


echo "[+] URL of the tarball: ${URL}"
echo "[+] Target name: ${NAME}"

echo "[+] Creating new server in rescue mode..."
SERVER=$(
    scw create \
        --bootscript=rescue \
        --name="[testing] $NAME" \
        --env="boot=rescue rescue_image=${URL}" \
        1GB
      )
echo "[+] Server created: ${SERVER}"

echo "[+] Booting..."
scw start "${SERVER}" >/dev/null
echo "[+] Done"
