FROM golang:1.5.3-alpine

ADD main.go /go/src/app/main.go

RUN go install app

ENTRYPOINT app
