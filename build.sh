#!/bin/bash

# Build all executables
#
# ex: ./build.sh
#
# 386 = 32 bit
# amd = 64 bit
#
# $GOOS			$GOARCH
# ====================
# darwin		386
# darwin		amd64
# linux			386
# linux			amd64
# windows		386
# windows		amd64

type setopt >/dev/null 2>&1 && setopt shwordsplit
PLATFORMS="darwin/386 darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64"

function go-compile {
	local GOOS=${1%/*}
	local GOARCH=${1#*/}
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o coreos-deploy-client-${GOOS}-${GOARCH} -i
}

function run {
	for PLATFORM in $PLATFORMS; do
			local CMD="go-compile ${PLATFORM}"
			echo "$CMD"
			$CMD
	done
}

run

mv coreos-deploy-client-windows-386 coreos-deploy-client-windows-386.exe
mv coreos-deploy-client-windows-amd64 coreos-deploy-client-windows-amd64.exe
