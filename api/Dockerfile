FROM golang:1.17-alpine

WORKDIR /api

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

WORKDIR /api/cmd/webserver

RUN go build -o main .

EXPOSE 8080

CMD ["/api/cmd/webserver/main"]