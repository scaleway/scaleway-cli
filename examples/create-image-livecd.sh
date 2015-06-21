#!/bin/bash

set -e
URL="${1}"

if [ -z "${1}" ]; then
    echo "usage: $(basename ${0}) <url>"
    echo ""
    echo "examples:"
    echo "  - $(basename ${0}) http://test-images.fr-1.storage.online.net/scw-distrib-ubuntu-trusty.tar"
    exit 1
fi

# FIXME: add usage

NAME=$(basename "${URL}")
NAME=${NAME%.*}

echo "[+] URL of the tarball: ${URL}" >&2
echo "[+] Target name: ${NAME}" >&2

echo "[+] Creating new server in live mode..." >&2
SERVER=$(
    scw create \
        --bootscript=stable \
        --name="[live] $NAME" \
        --env="boot=live rescue_image=${URL}" \
        50GB
      )
echo "[+] Server created: ${SERVER}" >&2

echo "[+] Booting..." >&2
scw start "${SERVER}" >/dev/null
echo "[+] Done" >&2

echo "${SERVER}"
