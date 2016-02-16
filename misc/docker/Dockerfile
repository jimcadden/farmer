FROM ubuntu:trusty

RUN apt-get update -qq
RUN apt-get install -yqq mercurial git gcc net-tools tcpdump iputils-ping golang

RUN mkdir -p /go/workspace/src
ENV PATH $PATH:/go/golang/bin:/go/workspace/bin:/go/until
ENV GOPATH /go/workspace

# build and install go circuit
RUN go get github.com/jmcadden/circuit/cmd/circuit
