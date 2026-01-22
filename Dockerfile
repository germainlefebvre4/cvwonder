FROM alpine:3 AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG JQ_VERSION=1.7
ARG CVWONDER_VERSION=0.3.0

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
        curl

# Install dependencies: >=jq-1.8.1
RUN apk add jq && \
    echo "jq version: " && \
    jq --version

# Download cvwonder binary for the target platform
RUN OS=$(echo $TARGETPLATFORM | cut -d'/' -f1) && \
    ARCH=$(echo $TARGETPLATFORM | cut -d'/' -f2) && \
    ARM_VERSION=$(echo $TARGETPLATFORM | cut -d'/' -f3) && \
    curl -L --output /app/cvwonder "https://github.com/germainlefebvre4/cvwonder/releases/download/v${CVWONDER_VERSION}/cvwonder_${OS}_${ARCH}" && \
    chmod +x /app/cvwonder && \
    /app/cvwonder version

FROM alpine:3

COPY --from=build /app/cvwonder /usr/local/bin/cvwonder

WORKDIR /cv

ENTRYPOINT ["cvwonder"]
CMD ["serve", "--input=cv.yml", "--output=generated/", "--theme=default", "--watch"]
