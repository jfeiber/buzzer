FROM golang:1.7

#Setup go stuff
RUN go get github.com/codegangsta/gin
RUN go get github.com/tools/godep
RUN go get bitbucket.org/liamstask/goose/cmd/goose

EXPOSE 3000
