docker pull ss251/message
docker rm -f message
docker rm -f mongodb

docker run -d \
-p 27017:27017 \
--name mongodb \
--network site \
mongo


docker run -d \
-e MONGOADDR="mongodb://mongodb:27017/chat" \
-e MYSQLADDR="userstore,root,password,users" \
--network site \
--name message \
ss251/message

#ssh -i "shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com