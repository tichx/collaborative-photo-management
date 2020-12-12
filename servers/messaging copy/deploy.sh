sh build.sh

docker push dayange/messaging

# ssh ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com "
ssh -i "c:\Users\Jesus\Downloads\info441win.pem" ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com "

docker rm -f mongodb
docker rm -f messaging
docker rm -f mq
docker pull dayange/messaging

export RABBIT_ADDR=mq:5672
export MQ=mq

docker run -d \
    -p 27017:27017 \
    --name mongodb \
    --network mynetwork \
    mongo

docker run -d \
    --name mq \
    --network mynetwork \
    -p 5672:5672 \
    -p 15672:15672 \
    rabbitmq:3-management

sleep 10

docker run -d \
    --name messaging \
    --network mynetwork \
    -e RABBIT_ADDR=$RABBIT_ADDR \
    -e MQ=$MQ \
    dayange/messaging
"