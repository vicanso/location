# location

[![Build Status](https://img.shields.io/travis/vicanso/location.svg?label=linux+build)](https://travis-ci.org/vicanso/location)


Get the location by ip address. IP data comes from [ip2region](https://github.com/lionsoul2014/ip2region).


```bash
curl 'http://127.0.0.1:7001/ip-location/json/1.0.132.192'
```

## start

```bash
docker run -d -p 7001:7001 vicanso/location
```