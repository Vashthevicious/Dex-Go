@ECHO OFF

set GOOS=linux
set GOARCH=amd64
go fmt ./...
go install -v