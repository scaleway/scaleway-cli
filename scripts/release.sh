#!/bin/bash

cd scripts/release || (echo "Please run this script from repo root" && exit 1)

yarn install --frozen-lock-file
yarn run release
