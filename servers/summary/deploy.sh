chmod +x build.sh
./build.sh
docker push ss251/summary
ssh -i "~/.ssh/shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com < 'server-deploy.sh'

#ssh -i "shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com