FROM golang:1.22

RUN apt-get update && apt-get upgrade -y

# install golangci lint
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1

# install node
RUN apt-get install -y npm
