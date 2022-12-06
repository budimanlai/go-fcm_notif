#!/bin/bash

GIT=`git log | head -n 1 | cut  -f 2 -d ' ' | head -c 8`

MAIN="*.go"
APP="fcm_notif"
FILE="dist/$APP-$GIT"

GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${FILE} $MAIN

upx ${FILE}

