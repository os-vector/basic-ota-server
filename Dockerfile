FROM golang:latest

WORKDIR /usr/src/basic-ota-server/

COPY . .
RUN go build -v -o /usr/local/bin/basic-ota-server ./...

CMD ["basic-ota-server"]

