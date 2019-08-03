FROM alpine:latest

# Install parallel and accept the citation notice (we aren't using this in a
# context where it make sense to cite GNU Parallel).
RUN apk add --no-cache parallel && \
    mkdir -p ~/.parallel && touch ~/.parallel/will-cite

RUN apk add bats
RUN apk add nfs-utils