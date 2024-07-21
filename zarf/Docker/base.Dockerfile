FROM golang:1.22
WORKDIR /tmp

RUN apt-get update && apt-get upgrade -y

# install requirements
RUN apt-get install -y clang gcc-aarch64-linux-gnu gcc-mingw-w64-x86-64 xz-utils gcc-arm-linux-gnueabi

## ADD ZIG to use as cross compiler for windows arm64
WORKDIR /opt
# download from: https://ziglang.org/download/
RUN wget https://ziglang.org/builds/zig-linux-x86_64-0.14.0-dev.367+a57479afc.tar.xz
RUN  tar -xJf zig-linux-x86_64-0.14.0-dev.367+a57479afc.tar.xz && \
     rm zig-linux-x86_64-0.14.0-dev.367+a57479afc.tar.xz && \
     mv zig-linux-x86_64-0.14.0-dev.367+a57479afc zig

ENV PATH="/opt/zig:${PATH}"

# install some utilities
RUN apt-get install -y joe bash-completion

# install golangci lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
## install go-licence-detectors
RUN go install go.elastic.co/go-licence-detector@latest

## install goreleaser oss
RUN wget https://github.com/goreleaser/goreleaser/releases/download/v2.1.0/goreleaser_2.1.0_amd64.deb && \
    dpkg -i goreleaser_2.1.0_amd64.deb
# install node
RUN apt-get install -y npm

WORKDIR /project