FROM golang:1.16.2-alpine3.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o code-swamp-server cmd/main.go
RUN go get github.com/mibk/dupl/...

FROM alpine:3.13
COPY --from=builder /go/bin/dupl /go/bin/dupl
COPY --from=builder /build /build

ENTRYPOINT  [ "./build/code-swamp-server" ]
