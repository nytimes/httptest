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
  -X main.BuildBranch=${DRONE_BRANCH} \
  -X main.BuildCommit=${DRONE_COMMIT:0:8} \
  -X main.BuildTime=$(date -Iseconds)"

# Distroless; smaller than Alpine, has SSL included, works for multi-arch
FROM gcr.io/distroless/static-debian12

# Copy binary from build container
COPY --from=build /go/bin/httptest /bin/httptest

# Default command
CMD ["/bin/httptest"]
