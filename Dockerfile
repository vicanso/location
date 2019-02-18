FROM node:10-alpine as webbuilder

RUN apk update \
  && apk add git \
  && git clone --depth=1 https://github.com/vicanso/location.git /location \
  && cd \location\web \
  && npm i \
  && npm run build

FROM golang:1.11-alpine as builder

COPY --from=webbuilder /location /location

RUN apk update \
  && apk add make \
  && go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /location \
  && packr2 \
  && make build

FROM alpine

EXPOSE 7001

COPY --from=builder /location/location /usr/local/bin/location

CMD ["location"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD location --mode=check || exit 1