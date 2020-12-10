GOOS=linux go build
docker build -t ss251/upload .
go clean

docker push ss251/upload

ssh ec2-user@ec2-3-211-194-149.compute-1.amazonaws.com < update.sh

printf "built upload"
exit