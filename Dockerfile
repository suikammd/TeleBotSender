FROM golang:1-alpine AS builder

RUN mkdir -p /go/src/github.com/suikammd/teleBotSender
COPY . /go/src/github.com/suikammd/teleBotSender

RUN apk upgrade \
    && apk add git \
    && go get github.com/suikammd/teleBotSender

FROM alpine:latest AS dist
LABEL maintainer="suikammd <suikammd@gmail.com>"

ENV BOTTOKEN b
ENV GROUPID 0
ENV PERSONALID 0

COPY --from=builder /go/bin/teleBotSender /usr/bin/teleBotSender

CMD exec teleBotSender \
        --b=$BOTTOKEN \
        --g=$GROUPID \
        --p=$PERSONALID