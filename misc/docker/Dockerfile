FROM ubuntu:trusty

RUN apt-get update -yqq
RUN apt-get install -yqq git gcc golang \
net-tools tcpdump iputils-ping \
qemu-kvm uml-utilities bridge-utils dnsmasq

RUN mkdir -p /go/workspace/src
ENV PATH $PATH:/go/golang/bin:/go/workspace/bin:/go/until
ENV GOPATH /go/workspace

# build and install go circuit
RUN go get github.com/jmcadden/circuit/cmd/circuit

#enable udp broadcast
RUN echo >> net.ipv4.icmp_echo_ignore_broadcasts = 0
RUN sysctl -p
