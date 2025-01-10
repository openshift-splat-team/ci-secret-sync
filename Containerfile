FROM registry.ci.openshift.org/openshift/release:golang-1.22 AS builder 

WORKDIR /go/src/github.com/openshift-splat-team/ci-secret-sync
COPY . .
RUN make build
