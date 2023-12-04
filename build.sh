tag=$1

docker build -f Dockerfile.production -t storm-anton:$tag .
docker image tag storm-anton:$tag 070998/storm-idx:$tag
docker push 070998/storm-idx:$tag
