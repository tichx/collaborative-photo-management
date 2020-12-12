GOOS=linux go build
docker build -t ss251/s3 .
go clean

docker push ss251/s3

ssh ec2-user@ec2-3-211-194-149.compute-1.amazonaws.com < update.sh

printf "built s3"
exit