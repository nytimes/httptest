# Build container
FROM golang:alpine

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY . /go/src/github.com/nytimes/httptest
WORKDIR /go/src/github.com/nytimes/httptest

# --build-arg
ARG DRONE_BRANCH
ARG DRONE_COMMIT

# Get Go packages
RUN apk add --no-cache git
RUN go get github.com/youmark/pkcs8

# Build application
RUN go build -a -o /go/bin/httptest \
    -ldflags "-extldflags \"-static\" \
              -X main.BuildBranch=${DRONE_BRANCH} \
              -X main.BuildCommit=${DRONE_COMMIT:0:8} \
              -X main.BuildTime=$(date -Iseconds)"

# Minimum runtime container
FROM alpine

# Install packages
RUN apk add --no-cache ca-certificates

# Copy built binary from build container
COPY --from=0 /go/bin/httptest /bin/httptest

# Default command
CMD ["/bin/httptest"]
