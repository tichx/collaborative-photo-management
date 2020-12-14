docker rm -f userstore
docker pull ss251/userstore

docker run -d \
-p 3306:3306 \
--name userstore \
-e MYSQL_ROOT_PASSWORD=password \
-e MYSQL_DATABASE=Users \
--network site \
ss251/userstore

#docker exec -t -i userstore /bin/bash -c "mysql -uroot -p$MYSQL_ROOT_PASSWORD"
#ssh -i "shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com


