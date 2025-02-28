FROM alpine:latest as build

ARG DISTRIBUTION=linux
ARG CPU_ARCH=amd64
ARG CVWONDER_VERSION=latest

WORKDIR /app

RUN apk update && \
    apk add --no-cache \
        wget curl
RUN wget https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
    mv jq-linux64 /usr/local/bin/jq && \
    chmod +x /usr/local/bin/jq

RUN VERSION=$(curl -s "https://api.github.com/repos/germainlefebvre4/cvwonder/releases/tags/${CVWONDER_VERSION}" | jq -r '.tag_name') && \
    curl -L -o /app/cvwonder "https://github.com/germainlefebvre4/cvwonder/releases/download/${VERSION}/cvwonder_${DISTRIBUTION}_${CPU_ARCH}" && \
    chmod +x cvwonder


FROM alpine:latest

COPY --from=build /app/cvwonder /usr/local/bin/cvwonder

WORKDIR /cv

ENTRYPOINT ["cvwonder"]
CMD ["serve", "--input=cv.yml", "--output=generated/", "--theme=default", "--watch"]
