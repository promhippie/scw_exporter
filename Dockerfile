FROM --platform=$BUILDPLATFORM golang:1.24.5-alpine3.21@sha256:6edc20586dd08dacad538c1f09984bc2aa61720be59056cf75429691f294d731 AS builder

RUN apk add --no-cache -U git curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/exporter
COPY . /go/src/exporter/

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 exporter && \
    adduser -D -u 1337 -h /var/lib/exporter -G exporter exporter

EXPOSE 9503
VOLUME ["/var/lib/exporter"]
ENTRYPOINT ["/usr/bin/scw_exporter"]
HEALTHCHECK CMD ["/usr/bin/scw_exporter", "health"]

COPY --from=builder /go/src/exporter/bin/scw_exporter /usr/bin/scw_exporter
WORKDIR /var/lib/exporter
USER exporter
