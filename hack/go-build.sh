#!/usr/bin/env bash

set -eu
go mod tidy
go mod vendor
go mod verify

REPO=github.com/openshift-splat-team/ci-secret-sync
WHAT=${1:-ci-secret-sync}
export GO111MODULE=on
#export GOFLAGS=${GOFLAGS:--mod=vendor}
GLDFLAGS=${GLDFLAGS:-}

eval $(go env | grep -e "GOHOSTOS" -e "GOHOSTARCH")

: "${GOOS:=${GOHOSTOS}}"
: "${GOARCH:=${GOHOSTARCH}}"

# Go to the root of the repo
cd "$(git rev-parse --show-cdup)"

if [ -z ${VERSION_OVERRIDE+a} ]; then
	if [ -n "${BUILD_VERSION+a}" ] && [ -n "${BUILD_RELEASE+a}" ]; then
		echo "Using version from the build system..."
		VERSION_OVERRIDE="${BUILD_VERSION}-${BUILD_RELEASE}"
	else
		echo "Using version from git..."
		VERSION_OVERRIDE=$(git describe --abbrev=8 --dirty --always)
	fi
fi

GLDFLAGS+="-extldflags '-static' -X ${REPO}/pkg/version.Raw=${VERSION_OVERRIDE}"

eval $(go env)

#export GOFLAGS=""
echo "Building ${REPO}/cmd/${WHAT} (${VERSION_OVERRIDE})"
#GO111MODULE=${GO111MODULE} CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOFLAGS} -ldflags "${GLDFLAGS}" -o bin/${WHAT} ${REPO}/cmd/${WHAT}
GO111MODULE=${GO111MODULE} CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build ${GOFLAGS} -ldflags "${GLDFLAGS}" -o bin/${WHAT} ./cmd/${WHAT}