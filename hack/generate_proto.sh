#!/usr/bin/env sh
set -o errexit
set -o nounset
set -o pipefail
protoc --proto_path=${GOPATH}/src:. --micro_out=pkg/ --go_out=pkg/ proto/*.proto
