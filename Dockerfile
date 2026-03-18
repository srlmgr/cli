# This file is used by goreleaser
FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/cli"]
COPY $TARGETPLATFORM/cli /
