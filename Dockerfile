# syntax=docker/dockerfile:1

# ---- Build stage: compile cvwonder from the (patched) sources ----------------
# Runs on the builder's native platform and cross-compiles, so multi-arch
# builds do not need to emulate the Go toolchain under QEMU.
FROM --platform=$BUILDPLATFORM golang:1.27.4-alpine AS build

ARG TARGETOS
ARG TARGETARCH
ARG CVWONDER_VERSION=dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath \
      -ldflags "-s -w -X github.com/germainlefebvre4/cvwonder/internal/version.CVWONDER_VERSION=v${CVWONDER_VERSION#v}" \
      -o /out/cvwonder ./cmd/cvwonder

# ---- Runtime stage: only chromium + the static binary ------------------------
FROM alpine:3

# chromium is required to render PDFs / screenshots. No curl, jq or build
# tooling is shipped in the final image.
RUN apk upgrade --no-cache && \
    apk add --no-cache chromium && \
    addgroup -g 1000 -S cvwonder && \
    adduser -u 1000 -S -G cvwonder -h /home/cvwonder cvwonder && \
    mkdir -p /cv && chown cvwonder:cvwonder /cv

COPY --from=build /out/cvwonder /usr/local/bin/cvwonder

USER cvwonder
WORKDIR /cv

ENV CHROME_BIN=/usr/bin/chromium-browser

ENTRYPOINT ["cvwonder"]
CMD ["serve", "--input=cv.yml", "--output=generated/", "--theme=default", "--watch"]
