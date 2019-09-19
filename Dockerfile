FROM node:10-alpine as webbuilder

RUN apk update \
  && apk add git \
  && git clone --depth=1 https://github.com/vicanso/location.git /location \
  && cd /location/web \
  && npm i \
  && npm run build \
  && rm -rf node_modules

FROM golang:1.13-alpine as builder

COPY --from=webbuilder /location /location

RUN apk update \
  && apk add git make \
  && go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /location \
  && make build

FROM alpine

EXPOSE 7001

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates

COPY --from=builder /location/location /usr/local/bin/location

USER go

WORKDIR /home/go

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 CMD [ "wget", "http://127.0.0.1:7001/ping", "-q", "-O", "-"]

CMD ["location"]
