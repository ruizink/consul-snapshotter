FROM scratch as scratch

ARG TARGETOS TARGETARCH

COPY --chmod=0755 build/bin/$TARGETOS/$TARGETARCH /usr/bin

ENTRYPOINT ["/usr/bin/consul-snapshotter"]

FROM alpine:latest as alpine

ARG TARGETOS TARGETARCH

COPY --chmod=0755 build/bin/$TARGETOS/$TARGETARCH /usr/bin

ENTRYPOINT ["/usr/bin/consul-snapshotter"]