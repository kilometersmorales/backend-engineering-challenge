FROM golang:1.21
WORKDIR /app/src/unbabel/
COPY . .
RUN go mod download
RUN go build -o unbabel .
RUN chmod +x unbabel
CMD ./unbabel --window_size 10
