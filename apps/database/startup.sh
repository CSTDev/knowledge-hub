docker kill database
docker rm database
docker run -d --name database \
	-p 27017:27107 \
	mongo:3.6
