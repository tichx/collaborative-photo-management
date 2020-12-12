# GOOS=linux go build
# docker build -t dayange/gateway .
# go clean

# $Env:GOOS = "linux"; go build .\main.go

go install
go build
$Env:GOOS = "linux";
docker build -t dayange/gateway .
go clean

docker push dayange/gateway
ssh -i "c:\Users\Jesus\Downloads\info441win.pem" ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com