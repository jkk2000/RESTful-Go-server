FROM golang:buster

RUN go get -u github.com/gorilla/mux

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /http-server

COPY . .

RUN go mod init github.com/jkk2000/RESTful-Go-server

RUN go build -o main

EXPOSE 8080

CMD ["./main"]