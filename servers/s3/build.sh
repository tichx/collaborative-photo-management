GOOS=linux go build
docker build -t ss251/s3 .
go clean

docker push ss251/s3

ssh -i "shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com < update.sh

printf "built s3"
exit