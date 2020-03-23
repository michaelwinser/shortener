FROM golang:1.14-alpine AS builder

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ADD . /src/
RUN cd /src && go install main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/main /go/bin/main

ENTRYPOINT ["/go/bin/main"]