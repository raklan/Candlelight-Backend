webcontid=$(docker ps -aqf "name=candlelight-backend-web-1")
rediscontid=$(docker ps -aqf "name=candlelight-backend-redis-1")
imageid=$(docker images candlelight-backend --format "{{.ID}}")

docker stop "$webcontid" && docker rm "$webcontid"
docker stop "$rediscontid" && docker rm "$rediscontid"
docker rmi "$imageid"

docker compose build
docker compose up