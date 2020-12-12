docker rm -f gateway
docker rm -f redis
docker pull tichx/gateway
docker pull tichx/message

docker run -d \
-p 6379:6379 \
--name redis \
--network site \
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
-e DSN='root:password@tcp(34.217.136.38:3306)/users' \
-e SESSIONKEY=arandomkeyforhashing \
-e REDISADDR=redis:6379 \
tichx/gateway

