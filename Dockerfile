# Build container
FROM --platform=$BUILDPLATFORM golang:alpine AS build

ARG TARGETOS TARGETARCH

ENV CGO_ENABLED=0

COPY . /go/src/github.com/nytimes/httptest
WORKDIR /go/src/github.com/nytimes/httptest

# --build-arg
ARG DRONE_BRANCH
ARG DRONE_COMMIT

# Build application
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o /go/bin/httptest \
  -ldflags "-extldflags \"-static\" \
  -X main.BuildBranch=${env.branch} \
  -X main.BuildCommit=${env.commit_sha} \
  -X main.BuildTime=$(date -Iseconds)"

# We can't use distroless because some teams need to add a bearer token to
# authenticate when the tests are run
FROM alpine

# Install dependencies
RUN apk add --no-cache ca-certificates

# Copy binary from build container
COPY --from=build /go/bin/httptest /bin/httptest

# Default command
CMD ["/bin/httptest"]
