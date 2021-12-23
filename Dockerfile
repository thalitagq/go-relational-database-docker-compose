FROM golang:latest
ENV GO111MODULE=off
RUN mkdir /app
WORKDIR /app

RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm