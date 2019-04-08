# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.
FROM golang:1.11-alpine3.7 as builder

# install tools
RUN apk add --no-cache git

WORKDIR /go/src/openpitrix.io/watcher
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN mkdir -p /openpitrix_bin
RUN go build -v -a -installsuffix cgo -ldflags '-w' -o /openpitrix_bin/watch pkg/watch/main.go


FROM alpine:3.7
COPY --from=builder /openpitrix_bin/watch /usr/local/bin/

CMD ["/usr/local/bin/watch"]
