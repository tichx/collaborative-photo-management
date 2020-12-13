GOOS=linux go build
docker build -t ss251/gateway .
go install
go clean