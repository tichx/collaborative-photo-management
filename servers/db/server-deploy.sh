docker rm -f userstore
docker pull ss251/userstore

docker run -d \
-p 3306:3306 \
--name userstore \
-e MYSQL_ROOT_PASSWORD=password \
-e MYSQL_DATABASE=users \
--network site \
ss251/userstore



