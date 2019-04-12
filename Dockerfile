FROM golang:alpine

RUN apk add --no-cache git mercurial

Run git clone https://github.com/golang/time /go/src/golang.org/x/time

LABEL max-Up-DownloadToOss.version="1.0.14" maintainer="Pharber"

ENV BM_HOME /go/bin

RUN go get github.com/alfredyang1986/blackmirror && \
go get github.com/alfredyang1986/BmServiceDef && \
go get github.com/PharbersDeveloper/max-Up-DownloadToOss

RUN go install -v github.com/PharbersDeveloper/max-Up-DownloadToOss

ADD resource /go/bin/resource
ADD tmp /go/bin/tmp

WORKDIR /go/bin

ENTRYPOINT ["max-Up-DownloadToOss"]
