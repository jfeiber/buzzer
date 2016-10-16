FROM golang:1.7

#Setup go stuff
RUN go get github.com/codegangsta/gin
RUN go get github.com/tools/godep
RUN go get bitbucket.org/liamstask/goose/cmd/goose

ENV SESSION_AUTHENTICATION_KEY hV3vC5AWX39IVUWSP2NcHciWvqZTa2N95RxRTZHWUsaD6HEdz0

EXPOSE 3000
