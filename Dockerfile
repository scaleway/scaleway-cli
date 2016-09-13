FROM golang:1.7
COPY . /go/src/github.com/scaleway/scaleway-cli
WORKDIR /go/src/github.com/scaleway/scaleway-cli
RUN go install -v ./cmd/scw
ENTRYPOINT ["scw"]
CMD ["help"]
