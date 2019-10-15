FROM golang:alpine AS build-env

# Set go compiler to use modules
ENV GO111MODULE=on

RUN apk update -qq && apk add git

WORKDIR $GOPATH/src/github.com/harvard-arts-reverse-proxy

COPY go.mod .

COPY go.sum .

RUN go mod tidy && go mod vendor

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

COPY . .

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o bin/harvard-arts-reverse-proxy" -command="./bin/harvard-arts-reverse-proxy"  -color -graceful-kill

EXPOSE 4000
