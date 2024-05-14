FROM golang:alpine

RUN apk update && apk upgrade && apk add aardvark-dns
WORKDIR /src
COPY . .

CMD go test -v ./...

