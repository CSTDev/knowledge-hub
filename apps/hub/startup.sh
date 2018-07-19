docker kill hub
docker rm hub
docker run -d \
	--name=hub \
	-p 3000:3000 \
	-e APP_CONFIG_VERSION="0.0.2" \
	-e APP_CONFIG_MAP_PROVIDER='http://a.tile.opencyclemap.org/cycle/${z}/${x}/${y}.png' \
	-e APP_CONFIG_API_URL='http://localhost:8000' \
	hub:0.0.3
