FROM golang:1.10.2-alpine3.7
RUN apk update; apk upgrade
RUN apk add git

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/ScaleSec/fedrampup
VOLUME ["/.aws"]

CMD "/go/bin/fedrampup"
