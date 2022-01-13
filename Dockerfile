
FROM golang:1.16.12-alpine3.14 AS builder
ENV GO111MODULE=on 
WORKDIR /app
COPY . ./
RUN go mod download

## Our project will now successfully build with the necessary go libraries included.
RUN go build -o storeapi .
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/storeapi .
EXPOSE 8080
CMD ["/app/storeapi"]