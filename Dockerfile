FROM golang:1.8 AS builder
WORKDIR /go/src/github.com/dlorch/nfsv3/
ADD ./ /go/src/github.com/dlorch/nfsv3/
RUN CGO_ENABLED=0 go build -o nfsv3-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/dlorch/nfsv3/nfsv3-server /usr/local/bin/
EXPOSE 111/udp
EXPOSE 111/tcp
EXPOSE 892/tcp
EXPOSE 2049/tcp
ENTRYPOINT ["/usr/local/bin/nfsv3-server"]