chmod +x build.sh
./build.sh
docker push tichx/summary
ssh -i "~/.ssh/shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com < 'server-deploy.sh'

