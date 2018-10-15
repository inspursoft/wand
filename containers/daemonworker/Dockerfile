FROM golang:1.9.6

MAINTAINER  wangkun_lc@inspur.com

COPY src/daemonworker /go/src/github.com/inspursoft/wand/src/daemonworker
COPY src/vendor /go/src/github.com/inspursoft/wand/src/vendor

WORKDIR /go/src/github.com/inspursoft/wand/src/daemonworker

RUN go build -v -a -o /go/bin/daemonworker && \
    chmod u+x /go/bin/daemonworker && \
    mkdir /data

WORKDIR /go/bin

VOLUME ["/data"]

CMD ["daemonworker"]

EXPOSE 8088


