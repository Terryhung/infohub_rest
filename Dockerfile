FROM golang:1.8

WORKDIR $GOPATH/src/github.com/Terryhung/infohub_rest
COPY . .

RUN go get -d -v
RUN go install -v

ENTRYPOINT /go/bin/infohub_rest

EXPOSE 8080
