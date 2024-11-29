FROM golang:alpine

LABEL maintainer = "github.com/h-dav"

RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base

RUN mkdir /app
WORKDIR /app

copy . .

RUN go get -d -v ./...

RUN go install -v ./...

# Enable hot reload
RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go get -v golang.org/x/tools/gopls

ENTRYPOINT CompileDaemon --build="go build -a -installsuffix cgo -o main ." --command=./main
