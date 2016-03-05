FROM golang:1.5.3-alpine

ENV HTTPS_PROXY http://web-proxy.austin.hpicorp.net:8080

ADD main.go /go/src/dockerfile/main.go

RUN go install dockerfile

ENTRYPOINT dockerfile
