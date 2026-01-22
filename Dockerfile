FROM alpine:3 AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG JQ_VERSION=1.7
ARG CVWONDER_VERSION=0.3.0

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
        curl

# Determine architecture for jq binary
RUN case "${TARGETPLATFORM}" in \
        "linux/amd64") JQ_ARCH="linux64" ;; \
        "linux/arm64") JQ_ARCH="linux64" ;; \
        "linux/arm/v7") JQ_ARCH="linux32" ;; \
        *) JQ_ARCH="linux64" ;; \
    esac && \
    curl --output jq https://github.com/stedolan/jq/releases/download/jq-${JQ_VERSION}/jq-${JQ_ARCH} && \
    mv jq /usr/local/bin/jq && \
    chmod +x /usr/local/bin/jq

# Determine architecture for cvwonder binary
RUN case "${TARGETPLATFORM}" in \
        "darwin") ARCH="darwin" ;; \
        "linux/amd64") ARCH="amd64" ;; \
        "linux/arm64") ARCH="arm64" ;; \
        "linux/arm/v7") ARCH="armv7" ;; \
        *) ARCH="amd64" ;; \
    esac && \
    VERSION=$(curl -s "https://api.github.com/repos/germainlefebvre4/cvwonder/releases/tags/v${CVWONDER_VERSION}" | jq -r '.tag_name') && \
    curl -L --output /app/cvwonder "https://github.com/germainlefebvre4/cvwonder/releases/download/${VERSION}/cvwonder_linux_${ARCH}" && \
    chmod +x /app/cvwonder


FROM alpine:3

COPY --from=build /app/cvwonder /usr/local/bin/cvwonder

WORKDIR /cv

ENTRYPOINT ["cvwonder"]
CMD ["serve", "--input=cv.yml", "--output=generated/", "--theme=default", "--watch"]
