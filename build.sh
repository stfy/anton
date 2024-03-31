tag=$1

docker buildx build --platform=linux/amd64 -f Dockerfile -t storm-anton:$tag .
docker image tag storm-anton:$tag ghcr.io/tsunami-exchange/anton:$tag
docker push ghcr.io/tsunami-exchange/anton:$tag
