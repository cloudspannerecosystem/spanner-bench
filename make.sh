#!/bin/bash

GOOS=windows GOARCH=amd64 go build -o "./bin/spannerbench_windows_amd64"
GOOS=linux GOARCH=amd64 go build -o "./bin/spannerbench_linux_amd64"
GOOS=darwin GOARCH=amd64 go build -o "./bin/spannerbench_darwin_amd64"
