FROM golang:1.16

LABEL version="0.1.0"

WORKDIR /go/src/github.com/ardafirdausr/discuss
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /go/bin/discuss cmd/discuss/*.go

ENTRYPOINT ["/go/bin/discuss"]