FROM golang:1.21.5-alpine3.17 AS builder

COPY . /test-infra
WORKDIR /test-infra
RUN go get github.com/onsi/ginkgo/v2/ginkgo/generators && \
    go install github.com/onsi/ginkgo/v2/ginkgo && \
    go get github.com/onsi/gomega/...

FROM golang:1.21.5-alpine3.17

COPY --from=builder /go/bin/ginkgo /ginkgo
RUN apk update && \
    apk add --no-cache ca-certificates git gcc musl-dev && \
    chmod +x /ginkgo
