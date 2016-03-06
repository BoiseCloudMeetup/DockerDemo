#!/bin/bash

pushd App
docker build -t jsparrow/chat-app .
popd

pushd KeyService
docker build -t jsparrow/key-service .
popd

pushd MessageService
docker build -t jsparrow/message-service .
popd

pushd UserService
docker build -t jsparrow/user-service .
popd

pushd WSService
docker build -t jsparrow/ws-service .
popd
