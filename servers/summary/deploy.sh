sh build.sh

docker push dayange/summary

ssh -i "c:\Users\Jesus\Downloads\info441win.pem" ec2-user@ec2-52-14-207-20.us-east-2.compute.amazonaws.com "

docker rm -f summary

docker pull dayange/summary

docker run -d \
--network mynetwork \
--name summary \
dayange/summary
"