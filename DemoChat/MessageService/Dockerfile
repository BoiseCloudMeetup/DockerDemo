FROM golang

ENV HTTPS_PROXY http://web-proxy.austin.hpicorp.net:8080

RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson

ADD main.go /go/src/app/main.go

RUN go install app

ENTRYPOINT app
