FROM golang:1.7

RUN go get github.com/codegangsta/gin

EXPOSE 3000

CMD ["gin"]
