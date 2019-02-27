# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM golang:1.11-alpine3.7 as builder

# intall tools
RUN apk add --no-cache git

WORKDIR /go/src/openpitrix.io/iam
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN mkdir -p /openpitrix_bin
RUN go generate openpitrix.io/iam/pkg/version && \
	GOBIN=/openpitrix_bin go install -ldflags '-w -s' -tags netgo openpitrix.io/iam/cmd/...

FROM alpine:3.7
COPY --from=builder /openpitrix_bin/* /usr/local/bin/
CMD ["/usr/local/bin/am"]
