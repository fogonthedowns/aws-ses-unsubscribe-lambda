#!/bin/bash

go get -u "github.com/fogonthedowns/aws-ses-unsubscribe-lambda/lib"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
zip main.zip main
rm main

