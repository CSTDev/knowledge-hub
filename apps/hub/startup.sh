docker kill hub
docker rm hub
docker run -d \
	--name=hub \
	-p 5000:5000 \
	-e REACT_APP_VERSION="0.0.2" \
	-e REACT_APP_MAP_PROVIDER='http://a.tile.opencyclemap.org/cycle/${z}/${x}/${y}.png' \
	-e REACT_APP_API_URL='http://localhost:9000' \
	hub:0.0.1
