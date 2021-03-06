FROM golang:alpine

RUN apk add --no-cache git mercurial

Run git clone https://github.com/golang/time /go/src/golang.org/x/time

LABEL  maintainer="developer@pharbers.com" max-Up-DownloadToOss.version="1.0.20"

ENV PH_FUAD_HOME $GOPATH/src/github.com/PharbersDeveloper/max-Up-DownloadToOss

RUN go get github.com/alfredyang1986/blackmirror && \
go get github.com/alfredyang1986/BmServiceDef && \
go get github.com/PharbersDeveloper/max-Up-DownloadToOss && \
go install -v github.com/PharbersDeveloper/max-Up-DownloadToOss

WORKDIR /go/bin

ENTRYPOINT ["max-Up-DownloadToOss"]
