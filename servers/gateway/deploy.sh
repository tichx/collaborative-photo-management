sh build.sh

docker push dayange/gateway

# ssh ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com "
ssh -i "c:\Users\Jesus\Downloads\info441win.pem" ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com "

docker rm -f gateway
docker rm -f redisServer
docker rm -f mysqldemo
docker network rm mynetwork
docker pull dayange/gateway
docker pull dayange/mysqldemo

export TLSCERT="/etc/letsencrypt/live/chatapi.danfengy.me/fullchain.pem" \
export TLSKEY="/etc/letsencrypt/live/chatapi.danfengy.me/privkey.pem" \
export MYSQL_ROOT_PASSWORD=123456
export MYSQL_ADDR=mysqldemo:3306
export REDISADDR=redisServer:6379
export SESSIONKEY=aKey
export SUMMARY_ADDR="summary:80"
export MESSAGE_ADDR="messaging:80"
export DSN="root:123456@tcp(mysqldemo:3306)/mysqldemo"
export NETWORK=mynetwork
export RABBIT_ADDR=mq:5672

docker network create mynetwork

docker run -d --name redisServer --network mynetwork redis

docker run -d \
    -p 3306:3306 \
        -e MYSQL_ROOT_PASSWORD=123456 \
    -e MYSQL_DATABASE=mysqldemo \
    --name mysqldemo \
    --network mynetwork \
    dayange/mysqldemo

sleep 10

docker run -d \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT="/etc/letsencrypt/live/chatapi.danfengy.me/fullchain.pem" \
-e TLSKEY="/etc/letsencrypt/live/chatapi.danfengy.me/privkey.pem" \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e SUMMARY_ADDR=$SUMMARY_ADDR \
-e MESSAGE_ADDR=$MESSAGE_ADDR \
-e REDISADDR=redisServer:6379 \
-e SESSIONKEY=$SESSIONKEY \
-e RABBIT_ADDR=$RABBIT_ADDR \
-e DSN=$DSN \
--network mynetwork \
--name gateway \
dayange/gateway

"