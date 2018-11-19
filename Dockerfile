FROM golang:1.11.2-alpine3.8
RUN apk update; apk upgrade
RUN apk add git

RUN go get -u github.com/golang/dep/cmd/dep
ADD . /go/src/fedrampup

WORKDIR /go/src/fedrampup
RUN /go/bin/dep ensure
RUN go install

CMD "/go/bin/fedrampup"
