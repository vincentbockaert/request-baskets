# Builds Request Baskets service using multi-stage builds
# Version 1.3

# Stage 1. Building service
FROM golang:latest as builder
WORKDIR /go/src/rbaskets
COPY . .
RUN GIT_VERSION="$(git describe --dirty='*' || git symbolic-ref -q --short HEAD)" \
  && GIT_COMMIT="$(git rev-parse HEAD)" \
  && GIT_COMMIT_SHORT="$(git rev-parse --short HEAD)" \
  && set -x \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
    -ldflags="-w -s -X main.GitVersion=${GIT_VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.GitCommitShort=${GIT_COMMIT_SHORT}" \
    -o /go/bin/rbaskets

# Stage 2. Packaging into distroless image to limit attack surface massively and reduce image size
FROM gcr.io/distroless/static-debian12

ENV PORT 8080
ENV LISTEN 0.0.0.0
# standard, adaptive, flatly
ENV THEME "flatly" 

VOLUME ["/var/lib/rbaskets"]
COPY --from=builder /go/bin/rbaskets /bin/rbaskets
EXPOSE ${PORT}
CMD ["/bin/rbaskets"]