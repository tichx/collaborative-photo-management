docker rm -f micro-upload

docker pull ss251/upload

docker run -d --name micro-upload --network 441network -p 8080:8080 ss251/upload