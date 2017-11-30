# Dockerfile.build represents the build environment of the driver, used during
# the development phase to test and in CI to build and test.

# The prefered base image is the lastest stable Alpine image, if alpine doesn't
# meet the requirements you can switch the from to the latest stable slim
# version of Debian (eg.: `debian:jessie-slim`).
FROM alpine:3.6

# To avoid files written in the volume by root or foreign users, we create a
# container local user with the same UID of the user executing the build.
# The following commands are defined to use in busybox based distributions,
# if you are using a standard distributions, replace the `adduser` command with:
#   `useradd --uid ${BUILD_UID} --home /opt/driver ${BUILD_USER}`
RUN mkdir -p /opt/driver/src && \
    adduser ${BUILD_USER} -u ${BUILD_UID} -D -h /opt/driver/src


# As minimal build tools you need: make, curl and git, install using the same
# command the specific tools required to build the driver.
RUN apk add --no-cache make git curl ca-certificates
RUN curl https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz -O && \
    tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/opt/driver"
ENV CONTAINER="yes"

# A nasty hack to make the Go binary work on alpine.
# See https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker.
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2    

# The volume with the full source code is mounted at `/opt/driver/src` so, we
# set the workdir to this path.
WORKDIR /opt/driver/src