# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.12-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/anonuuid
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM            alpine:3.10
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="anonuuid" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/anonuuid/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/anonuuid" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/anonuuid" \
                org.label-schema.help="docker exec -it $CONTAINER anonuuid --help"
COPY            --from=builder /go/bin/anonuuid /bin/
ENTRYPOINT      ["/bin/anonuuid"]
#CMD             []
