docker build -t ss251/photo-collab_db .

docker push ss251/photo-collab_db

ssh ec2-user@ec2-3-211-194-149.compute-1.amazonaws.com < update.sh

printf "built db"
exit