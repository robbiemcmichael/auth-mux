FROM golang:1.12.5-alpine3.9 AS builder

RUN apk update && apk add --no-cache git

ENV GO111MODULE=on
WORKDIR $GOPATH/src/github.com/robbiemcmichael/auth-mux

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-mux ./cmd/auth-mux

####

FROM alpine:3.9

COPY --from=builder /go/src/github.com/robbiemcmichael/auth-mux/auth-mux /usr/local/bin/auth-mux
