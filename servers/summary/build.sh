GOOS=linux go build
docker build -t ss251/summary .
go install
go clean