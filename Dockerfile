FROM golang:1.7

RUN go get github.com/codegangsta/gin
RUN go get github.com/tools/godep

EXPOSE 3000
