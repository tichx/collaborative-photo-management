docker pull ss251/photo
docker rm -f photo
# docker rm -f mongodb

# docker run -d \
# -p 27017:27017 \
# --name mongodb \
# --network site \
# mongo

docker run -d \
-e MONGOADDR="mongodb://mongodb:27017/photo" \
-e MYSQLADDR="userstore,root,password,users" \
--network site \
--name photo \
ss251/photo