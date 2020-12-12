docker pull tichx/summary
docker rm -f summary
docker run -d \
--network site \
--name summary \
tichx/summary