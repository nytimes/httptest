# Build container
FROM golang:alpine

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . /go/src/github.com/blupig/httptest
WORKDIR /go/src/github.com/blupig/httptest

# --build-arg
ARG BUILD_BRANCH
ARG BUILD_COMMIT

# Build application
RUN go build -a -o /go/bin/httptest \
    -ldflags "-extldflags \"-static\" \
              -X github.com/blupig/httptest/main.BuildBranch=${BUILD_BRANCH} \
              -X github.com/blupig/httptest/main.BuildCommit=${BUILD_COMMIT:0:8} \
              -X github.com/blupig/httptest/main.BuildTime=$(date -Iseconds)"

# Minimum runtime container
FROM alpine

# Install packages
RUN apk add --no-cache ca-certificates

# Copy built binary from build container
COPY --from=0 /go/bin/httptest /bin/httptest

# Default command
CMD ["/bin/httptest"]
