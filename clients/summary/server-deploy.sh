#!/usr/bin/env bash
docker rm -f client
docker pull ss251/client
docker run -d \
-p 80:80 -p 443:443 \
--name client \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-v /usr/share/nginx/ \
-e TLSKEY=/etc/letsencrypt/live/xutiancheng.me/privkey.pem \
-e TLSCERT=/etc/letsencrypt/live/xutiancheng.me/fullchain.pem \
ss251/client