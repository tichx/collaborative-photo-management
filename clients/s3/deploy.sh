chmod +x server-deploy.sh
chmod +x build.sh
./build.sh
docker push ss251/client

ssh -i "~/.ssh/shiny.pem" ubuntu@ec2-34-213-29-25.us-west-2.compute.amazonaws.com < 'server-deploy.sh'
