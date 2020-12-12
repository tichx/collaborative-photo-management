GOOS=linux go build
docker build -t tichx/summary .
go install
go clean