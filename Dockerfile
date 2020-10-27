FROM golang:1.15-alpine as builder

ENV BUILD_IN_DOCKER true

# ca-certificates is needed to add the certificates on the next image
# since it's FROM scratch, it does not have any certificates
# bash is needed to run the build script
RUN apk update && apk add --no-cache bash git

WORKDIR /go/src/github.com/scaleway/scaleway-cli

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY scripts/ scripts/
COPY cmd/ cmd/
COPY internal/ internal/
COPY .git/ .git/

RUN ./scripts/build.sh

FROM alpine:3.12
WORKDIR /
RUN apk update && apk add --no-cache bash ca-certificates openssh-client && update-ca-certificates
COPY --from=builder /go/src/github.com/scaleway/scaleway-cli/scw .
ENTRYPOINT ["/scw"]
