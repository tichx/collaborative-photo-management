# GOOS=linux go build
# docker build -t dayange/gateway .
# go clean

# $Env:GOOS = "linux"; go build .\main.go
$Env:GOOS = "linux";
docker build -t dayange/summary .
go clean