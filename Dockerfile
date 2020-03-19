FROM golang:1.14-alpine as builder

# ca-certificates is needed to add the certificates on the next image
# since it's FROM scratch, it does not have any certificates
# bash is needed to run the build script
RUN apk update && apk add --no-cache bash git ca-certificates && update-ca-certificates

WORKDIR /go/src/github.com/scaleway/scaleway-cli

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY scripts/ scripts/
COPY cmd/ cmd/
COPY internal/ internal/
COPY .git/ .git/

RUN bash scripts/build.sh build-in-docker

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/scaleway/scaleway-cli/scw .
ENTRYPOINT ["/scw"]
