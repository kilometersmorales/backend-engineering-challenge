FROM golang:1.21
WORKDIR /app/src/unbabel/
COPY . .
RUN go mod download
RUN go build -o unbabel .
RUN chmod +x unbabel
#CMD ./unbabelapp
ENTRYPOINT ["./unbabel", "-d", "--window_size", "10"]
