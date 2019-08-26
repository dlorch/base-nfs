FROM golang:1.8 AS builder
WORKDIR /go/src/github.com/dlorch/base-nfs/
ADD ./ /go/src/github.com/dlorch/base-nfs/
# CGO_ENABLED=0 necessary for the binary to run in alpine
RUN CGO_ENABLED=0 go build -o nfsv3-server

FROM alpine:latest
COPY --from=builder /go/src/github.com/dlorch/base-nfs/nfsv3-server /usr/local/bin/
EXPOSE 111/udp
EXPOSE 111/tcp
EXPOSE 892/tcp
EXPOSE 2049/tcp
ENTRYPOINT ["/usr/local/bin/nfsv3-server"]