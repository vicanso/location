FROM golang:1.11-alpine as builder

RUN apk update \
  && apk add git make gcc \
  && git clone --depth=1 https://github.com/vicanso/location.git /location \
  && cd /location \
  && make build

FROM alpine

EXPOSE 7001

COPY --from=builder /location/location /usr/local/bin/location

CMD ["location"]

HEALTHCHECK --interval=10s --timeout=3s \
  CMD location --mode=check || exit 1