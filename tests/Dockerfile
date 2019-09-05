# need older version of alpine, in order to have older version of nfs-utils
# which is compatible with NFSv3. Commands like 'showmount' don't offer any
# possibility of selecting the NFS version and would default to NFSv4, which
# is incompatible with base-nfs.
FROM alpine:3.7

RUN apk --no-cache add bats ncurses tcpdump nfs-utils
