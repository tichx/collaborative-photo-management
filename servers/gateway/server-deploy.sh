docker rm -f gateway
docker rm -f redis
docker pull ss251/gateway

docker run -d \
-p 6379:6379 \
--name redis \
--network site \
--sysctl net.core.somaxconn=511 \
redis 

docker run -d \
-p 443:443 \
--name gateway \
--network site \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSKEY=/etc/letsencrypt/live/api.xutiancheng.me/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/api.xutiancheng.me/fullchain.pem \
-e MESSAGESADDR=message:80 \
-e SUMMARYADDR=summary:80 \
-e PHOTOSADDR=photo:80 \
-e S3ADDR=micro-s3:8080 \
-e DSN='root:password@tcp(34.217.136.38:3306)/Users' \
-e SESSIONKEY=arandomkeyforhashing \
-e REDISADDR=redis:6379 \
ss251/gateway

#docker exec -t -i userstore /bin/bash -c "mysql -uroot -p$MYSQL_ROOT_PASSWORD"
#ssh -i "shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com

