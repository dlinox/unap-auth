set GOOS=linux
set GOARCH=amd64
go build cmd/main.go

scp C:\LINOX\GO\unap.auth.go\main.zip linox@38.43.133.27:/home/linox/
