FROM golang:1.23-alpine3.19 as builder
ENV GOPATH=/go/src
WORKDIR /go/src

COPY src network-bug
WORKDIR /go/src/network-bug
RUN go build -o network-bug

FROM alpine:latest

COPY --from=builder /go/src/network-bug/network-bug /app/network-bug

CMD ["/app/network-bug"]
