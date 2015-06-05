FROM golang:cross

ENV CGO_ENABLED 0

# Install Godep for vendoring
RUN go get github.com/tools/godep
# Recompile the standard library without CGO
RUN go install -a std

RUN apt-get install -y -q git

#RUN go get -d github.com/docker/docker github.com/Sirupsen/logrus github.com/dustin/go-humanize github.com/kardianos/osext golang.org/x/crypto/ssh/terminal

# Declare the maintainer
MAINTAINER Scaleway Team <opensource@scaleway.com> (@scaleway)

# For convenience, set an env variable with the path of the code
ENV APP_DIR /go/src/github.com/scaleway/scaleway-cli

ADD . /go/src/github.com/scaleway/scaleway-cli

#RUN cd $APP_DIR && godep get && godep save

# Compile the binary and statically link
RUN cd $APP_DIR && GOOS=darwin GOARCH=amd64          godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Darwin-x86_64
RUN cd $APP_DIR && GOOS=darwin GOARCH=386            godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Darwin-i386
RUN cd $APP_DIR && GOOS=linux  GOARCH=amd64          godep go build -a -v -ldflags '-w -s' -o /go/bin/scw-Linux-x86_64
RUN cd $APP_DIR && GOOS=linux  GOARCH=386            godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-i386
RUN cd $APP_DIR && GOOS=linux  GOARCH=arm   GOARM=5  godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-arm
RUN cd $APP_DIR && GOOS=linux  GOARCH=arm   GOARM=6  godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-armv6
RUN cd $APP_DIR && GOOS=linux  GOARCH=arm   GOARM=7  godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Linux-armv7
RUN cd $APP_DIR && GOOS=freebsd  GOARCH=amd64          godep go build -a -v -ldflags '-w -s' -o /go/bin/scw-Freebsd-x86_64
RUN cd $APP_DIR && GOOS=freebsd  GOARCH=386            godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Freebsd-i386
RUN cd $APP_DIR && GOOS=freebsd  GOARCH=arm   GOARM=5  godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Freebsd-arm
#RUN cd $APP_DIR && GOOS=openbsd  GOARCH=amd64          godep go build -a -v -ldflags '-w -s' -o /go/bin/scw-Openbsd-x86_64
#RUN cd $APP_DIR && GOOS=openbsd  GOARCH=386            godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Openbsd-i386
#RUN cd $APP_DIR && GOOS=openbsd  GOARCH=arm   GOARM=5  godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Openbsd-arm
RUN cd $APP_DIR && GOOS=windows  GOARCH=amd64          godep go build -a -v -ldflags '-w -s' -o /go/bin/scw-Windows-x86_64
#RUN cd $APP_DIR && GOOS=windows  GOARCH=386            godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Windows-i386
#RUN cd $APP_DIR && GOOS=windows  GOARCH=arm   GOARM=5  godep go build -a -v -ldflags '-d -w -s' -o /go/bin/scw-Windows-arm
