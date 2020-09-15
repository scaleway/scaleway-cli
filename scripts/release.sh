#!/bin/bash

export SCW_SECRET_KEY="$(vault kv get -field scaleway_secret_key agw_kv/prd)"
export SCW_ACCESS_KEY="$(vault kv get -field scaleway_access_key agw_kv/prd)"
export GITHUB_TOKEN="$(vault kv get -field cli_release_github_token front_kv/opensource)"
vault kv get -field cli_release_dockerhub_token front_kv/opensource | docker login --password-stdin -u "$(vault kv get -field dockerhub_bot_username front_kv/opensource)"
cd scripts/release || (echo "Please run this script from repo root" && exit 1)

yarn install --frozen-lock-file
yarn run release
