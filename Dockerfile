FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/scaleway-cli /usr/bin/scw
ENTRYPOINT ["/usr/bin/scw"]
