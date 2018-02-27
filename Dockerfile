FROM golang:1.8

WORKDIR /Users/terry/GoProjects/src/github.com/Terryhung/infohub_rest
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["infohub_rest"]
