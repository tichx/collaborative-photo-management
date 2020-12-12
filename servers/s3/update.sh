docker rm -f micro-s3

docker pull ss251/s3

docker run -d --name micro-s3 --network 441network -p 8080:8080 ss251/s3