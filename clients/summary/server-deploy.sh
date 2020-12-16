#!/usr/bin/env bash
docker rm -f client
docker pull dayange/client
docker run -d \
-p 80:80 -p 443:443 \
--name client \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-v /usr/share/nginx/ \
-e TLSKEY=/etc/letsencrypt/live/xutiancheng.me/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/xutiancheng.me/fullchain.pem \
dayange/client

#ssh -i "shiny.pem" ubuntu@ec2-34-213-29-25.us-west-2.compute.amazonaws.com