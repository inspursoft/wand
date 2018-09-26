FROM golang:1.9.6

MAINTAINER  wangkun_lc@inspur.com

COPY src/daemonworker /go/src/devopssuite/src/daemonworker
COPY src/vendor /go/src/devopssuite/vendor

WORKDIR /go/src/devopssuite/src/daemonworker

RUN go build -v -a -o /go/bin/daemonworker && \
    chmod u+x /go/bin/daemonworker && \
    mkdir /data

WORKDIR /go/bin

VOLUME ["/data"]

CMD ["daemonworker"]

EXPOSE 8088


