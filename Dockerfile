FROM golang:1.21
WORKDIR /app/src/unbabel/
COPY . .
RUN go mod download
RUN go build -o unbabelapp .
RUN chmod +x unbabelapp
CMD ./unbabelapp
