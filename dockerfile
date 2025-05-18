### BUILD STAGE ###
FROM golang:1.20.14-alpine3.18 AS builder

LABEL maintainer="Arda <ardafirdausr@gmail.com>"

WORKDIR /app

COPY . .
RUN go install -v ./...
RUN go build -o /app cmd/discuss/*.go

### RUN STAGE ###
FROM alpine:3.18

RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
