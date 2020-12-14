docker rm -f micro-s3

docker pull ss251/s3

docker run -d --name micro-s3 --network site ss251/s3 \

#Access Key ID:AKIAJOMUO3S2R36HG3JQ Secret Access Key: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
#ssh -i "shiny.pem" ubuntu@ec2-34-217-136-38.us-west-2.compute.amazonaws.com