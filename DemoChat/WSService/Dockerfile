FROM golang

ENV HTTPS_PROXY http://web-proxy.austin.hpicorp.net:8080

RUN go get github.com/gorilla/websocket

ADD main.go /go/src/ws/main.go
ADD connection.go /go/src/ws/connection.go

RUN go install ws

ENTRYPOINT  ws
