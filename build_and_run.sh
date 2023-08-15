docker stop unbabeljg
docker rm unbabeljg
docker build . -t unbabeljg
docker run -d --name unbabeljg unbabeljg:latest
docker logs -f unbabeljg
