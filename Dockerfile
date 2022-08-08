# This Dockerfile builds lbcwallet from source and creates a small (55 MB) docker container based on alpine linux.
#
# Clone this repository and run the following command to build and tag a fresh lbcwallet amd64 container:
#
# docker build . -t yourregistry/lbcwallet
#
# You can use the following command to buid an arm64v8 container:
#
# docker build . -t yourregistry/lbcwallet --build-arg ARCH=arm64v8
#
# For more information how to use this docker image visit:
# https://github.com/lbryio/lbcwallet/tree/master/docs
#
# 9244  Mainet Wallet RPC port

ARG ARCH=amd64

FROM golang:1.19 AS build-container

ARG ARCH

ADD . /app
WORKDIR /app
RUN set -ex \
  && if [ "${ARCH}" = "amd64" ]; then export GOARCH=amd64; fi \
  && if [ "${ARCH}" = "arm32v7" ]; then export GOARCH=arm; fi \
  && if [ "${ARCH}" = "arm64v8" ]; then export GOARCH=arm64; fi \
  && echo "Compiling for $GOARCH" \
  && go install -v . ./cmd/...

FROM $ARCH/debian:bullseye-20220418-slim

COPY --from=build-container /go/bin /bin

VOLUME ["/root/.lbcwallet"]

EXPOSE 9244

ENTRYPOINT ["lbcwallet"]
