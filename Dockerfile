FROM alpine

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG BUILDPLATFORM

RUN echo "BUILDPLATFORM=$BUILDPLATFORM" \
 && echo "TARGETPLATFORM=$TARGETPLATFORM" \
 && echo "TARGETOS=$TARGETOS" \
 && echo "TARGETARCH=$TARGETARCH"

RUN apk add --no-cache bash

# See https://docs.github.com/en/actions/reference/workflows-and-actions/dockerfile-support#workdir
# WORKDIR /app

COPY --chmod=755 \
${TARGETPLATFORM}/yaml-reader* \
gh-action-entrypoint.sh \
/app/

ENTRYPOINT ["/app/yaml-reader"]
