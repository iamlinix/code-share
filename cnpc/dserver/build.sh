#/bin/sh
echo ">>> building linux platform"
go build dserver.go
echo ">>> building windows platform"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build dserver.go
ls -al dserver*
