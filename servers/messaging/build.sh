
npm build

docker build -t dayange/messaging .

docker push dayange/messaging

ssh -i "c:\Users\Jesus\Downloads\info441win.pem" ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com