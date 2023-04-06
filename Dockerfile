FROM golang:alpine

WORKDIR /usr/src/app

ADD ["static", "static"]
ADD ["gochan", "gochan"]

RUN export CGO_CFLAGS_ALLOW='-Xpreprocessor' \
    && go mod init gochan

WORKDIR /usr/src/app/gochan
RUN go build

CMD ["./gochan"]
