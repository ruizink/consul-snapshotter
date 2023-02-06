FROM scratch

ARG TARGETOS TARGETARCH

COPY build/bin/$TARGETOS/$TARGETARCH/consul-snapshotter /

ENTRYPOINT ["/consul-snapshotter"]