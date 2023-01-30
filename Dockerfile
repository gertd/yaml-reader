ARG GO_VERSION
FROM ghcr.io/aserto-dev/go-builder:$GO_VERSION AS build

WORKDIR /src

# dowload debugger into Docker cacheable layer
ENV GOBIN=/bin
ENV ROOT_DIR=/src

# generate & build
ARG VERSION
ARG COMMIT
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=ssh \
    mage build
        
FROM alpine
ARG VERSION
ARG COMMIT
LABEL org.opencontainers.image.version=$VERSION
LABEL org.opencontainers.image.source=https://github.com/gertd/yaml-reader
LABEL org.opencontainers.image.title="yaml file reader"
LABEL org.opencontainers.image.revision=$COMMIT
LABEL org.opencontainers.image.url=https://ghcr/gertd/yaml-reader

RUN apk add --no-cache bash git openssh
WORKDIR /app
COPY --from=build /src/dist/build_linux_amd64_v1/yaml-reader /app/
COPY --from=build /src/gh-action-entrypoint.sh /app/

ENTRYPOINT ["/app/yaml-reader"]
