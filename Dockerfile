FROM golang:1.8
WORKDIR /go/src/github.com/dlorch/nfsv3/
ADD ./ /go/src/github.com/dlorch/nfsv3/
RUN go build -o nfsv3-server && mv ./nfsv3-server /go/bin/
EXPOSE 111/udp
EXPOSE 111/tcp
EXPOSE 892/tcp
EXPOSE 2049/tcp
ENTRYPOINT ["/go/bin/nfsv3-server"]
