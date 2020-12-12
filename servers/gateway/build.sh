GOOS=linux go build
docker build -t tichx/gateway .
go install
go clean