FROM scratch as scratch

ARG TARGETOS TARGETARCH

COPY --chmod=0755 build/bin/consul-snapshotter_${TARGETOS}_${TARGETARCH} /usr/bin/consul-snapshotter

ENTRYPOINT ["/usr/bin/consul-snapshotter"]

FROM alpine:latest as alpine

ARG TARGETOS TARGETARCH

COPY --chmod=0755 build/bin/consul-snapshotter_${TARGETOS}_${TARGETARCH} /usr/bin/consul-snapshotter

ENTRYPOINT ["/usr/bin/consul-snapshotter"]