FROM golang:1.14-alpine as builder

# ca-certificates is needed to add the certificates on the next image
# since it's FROM scratch, it does not have any certificates
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /go/src/github.com/scaleway/scaleway-cli

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY .git/ .git/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-w -extldflags -static -X main.GitCommit=$(git rev-parse --short HEAD) -X main.GitBranch=$(git symbolic-ref -q --short HEAD || echo HEAD) -X main.BuildDate=$(date -u '+%Y-%m-%dT%I:%M:%S%p')" ./cmd/scw

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/scaleway/scaleway-cli/scw .
ENTRYPOINT ["/scw"]
