docker kill knowledge
docker rm knowledge
docker run -d --name=knowledge \
    -p 8000:8000 \
    -e MONGODB_URI=http://localhost:27017 \
    knowledge:0.0.9