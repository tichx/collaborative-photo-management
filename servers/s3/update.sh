docker rm -f micro-s3

docker pull ss251/s3

docker run -d --name micro-s3 --network site -p 8081:8081 ss251/s3 \