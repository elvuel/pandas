#!/bin/bash
set -e -o pipefail

: "${WORKDIR:=./things/}"
: "${NOCOMPRESS:=false}"


GO_BINDATA="pushd ${WORKDIR} && \
                go-bindata-assetfs -pkg things -nocompress=${NOCOMPRESS} dist/... && \
                popd"

bash -c "${GO_BINDATA}"
