docker pull ss251/photo-colalb
docker rm -f PhotoCollab_DB
export MYSQL_ROOT_PASSWORD="sailesh123"
docker run -d -p 3306:3306 --name PhotoCollab_DB -e MYSQL_ROOT_PASSWORD="sailesh123" --network 441network -e MYSQL_DATABASE=PhotoCollab ss251/photo-collab