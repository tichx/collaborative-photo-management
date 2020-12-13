FROM nginx
COPY build /usr/share/nginx/html
ADD default.conf /etc/nginx/conf.d/default.conf
EXPOSE 443
EXPOSE 80

#docker build -t ss251/s3client .
#docker push ss251/s3client

#ssh ec2-user@ec2-54-160-41-239.compute-1.amazonaws.com

#docker rm -f s3
#docker pull ss251/s3client
#docker run -d --name s3 -p 80:80 -p 443:443 \
#-v /etc/letsencrypt:/etc/letsencrypt:ro ss251/s3client