FROM golang:1.10.2-alpine3.7
RUN apk update; apk upgrade
RUN apk add git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/app
VOLUME ["/go/src/app"]
VOLUME ["/.aws"]

RUN go install

# docker run -v "$(pwd)":/go/src/app
ENTRYPOINT /go/bin/fedrampup
