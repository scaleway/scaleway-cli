FROM golang:cross

ENV CGO_ENABLED 0

# Recompile the standard library without CGO
RUN go install -a std

RUN apt-get install -y -q git

# Declare the maintainer
MAINTAINER Scaleway Team <opensource@scaleway.com> (@scaleway)

# For convenience, set an env variable with the path of the code
ENV APP_DIR /go/src/github.com/scaleway/scaleway-cli
WORKDIR $APP_DIR

ADD . /go/src/github.com/scaleway/scaleway-cli


# Compile the binary and statically link
RUN  GOOS=darwin   GOARCH=amd64          go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Darwin-x86_64
RUN  GOOS=darwin   GOARCH=386            go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Darwin-i386
RUN  GOOS=linux    GOARCH=386            go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-i386
#RUN  GOOS=linux    GOARCH=amd64          go build -a -v -ldflags '-w -s'    -o /go/bin/scw-Linux-x86_64
RUN cp /go/bin/scw-Linux-i386 /go/bin/scw-Linux-x86_64
RUN  GOOS=linux    GOARCH=arm   GOARM=5  go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-arm
RUN  GOOS=linux    GOARCH=arm   GOARM=6  go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-armv6
RUN  GOOS=linux    GOARCH=arm   GOARM=7  go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-armv7
RUN  GOOS=freebsd  GOARCH=amd64          go build -a -v -ldflags '-w -s'    -o /go/bin/scw-Freebsd-x86_64
RUN  GOOS=freebsd  GOARCH=386            go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Freebsd-i386
RUN  GOOS=freebsd  GOARCH=arm   GOARM=5  go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Freebsd-arm
#RUN GOOS=openbsd  GOARCH=amd64          go build -a -v -ldflags '-w -s'    -o /go/bin/scw-Openbsd-x86_64
#RUN GOOS=openbsd  GOARCH=386            go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Openbsd-i386
#RUN GOOS=openbsd  GOARCH=arm   GOARM=5  go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Openbsd-arm
RUN  GOOS=windows  GOARCH=amd64          go build -a -v -ldflags '-w -s'    -o /go/bin/scw-Windows-x86_64.exe
#RUN GOOS=windows  GOARCH=386            go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Windows-i386
#RUN GOOS=windows  GOARCH=arm   GOARM=5  go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Windows-arm
