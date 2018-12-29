# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

# -----------------------------------------------------------------------------
# builder
# -----------------------------------------------------------------------------

FROM golang:1.11-alpine3.7 as builder

# intall tools
RUN apk add --no-cache git

WORKDIR /build-dir
COPY . .

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOBIN=/build-dir

RUN echo module drone > /build-dir/go.mod
RUN git describe --tags --always > /build-dir/version
RUN git describe --exact-match 2>/dev/null || git log -1 --format="%H" > /build-dir/version

RUN go get -ldflags '-w -s' -tags netgo openpitrix.io/iam@$(cat /build-dir/version)

RUN echo version: $(cat /build-dir/version)

# -----------------------------------------------------------------------------
# for image
# -----------------------------------------------------------------------------

FROM alpine:3.7

COPY --from=builder /build-dir/iam /usr/local/bin/

CMD ["sh"]

# -----------------------------------------------------------------------------
# END
# -----------------------------------------------------------------------------
