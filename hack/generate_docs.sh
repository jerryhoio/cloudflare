#!/usr/bin/env sh
set -o errexit
set -o nounset
set -o pipefail

rm -rf docs/*
hugo -s hugo/