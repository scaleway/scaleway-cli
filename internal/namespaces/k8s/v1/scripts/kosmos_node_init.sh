#!/usr/bin/env sh

set -e
echo "Downloading node agent binary..."
wget https://scwcontainermulticloud.s3.fr-par.scw.cloud/node-agent_linux_amd64 -q
echo "Success"
chmod +x node-agent_linux_amd64
export POOL_ID=<pool-id>  POOL_REGION=<pool-region>  SCW_SECRET_KEY=<secret-key>
sudo -E ./node-agent_linux_amd64 -loglevel 0 -no-controller