FROM scratch
MAINTAINER Scaleway <opensource@scaleway.com> (@scaleway)

# ssl.tar is created by `make prepare-release`. It is an extract of the
# /etc/ssl directory of a golang image.
ADD ssl.tar /

COPY latest/scw-linux-i386 /bin/scw

ENTRYPOINT ["/bin/scw"]
