FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/scw /usr/bin/scw
ENTRYPOINT ["/usr/bin/scw"]
