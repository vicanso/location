package controller

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"

	"github.com/vicanso/cod"
	"github.com/vicanso/location/router"
	"github.com/vicanso/location/service"
)

type (
	// ipLocationCtrl ip location ctrl
	ipLocationCtrl struct{}
)

var (
	unknownIPLocation = &service.IPLocation{}
)

func init() {
	g := router.NewGroup("/ip-location")
	ctrl := ipLocationCtrl{}
	g.GET("/:type/:ip", ctrl.getLocation)
}

func (ctrl ipLocationCtrl) getLocation(c *cod.Context) (err error) {
	c.NoCache()
	ipAddr := c.Param("ip")
	// 如果为此地址，则使用客户端IP
	if ipAddr == "127.0.0.1" {
		ipAddr = c.RealIP()
	}
	ip, err := service.ConvertIP2Uint32(ipAddr)
	if err != nil {
		return
	}
	location := service.GetLocationByIP(ip)
	if location == nil {
		location = unknownIPLocation
	}
	switch c.Param("type") {
	case "xml":
		buf, err := xml.Marshal(location)
		if err != nil {
			return err
		}
		c.SetHeader(cod.HeaderContentType, "text/xml; charset=UTF-8")
		c.BodyBuffer = bytes.NewBuffer(buf)
	case "jsonp":
		callback := c.QueryParam("callback")
		if callback == "" {
			return errors.New("callback can not be null")
		}
		buf, err := json.Marshal(location)
		if err != nil {
			return err
		}
		c.SetHeader(cod.HeaderContentType, "application/javascript; charset=UTF-8")
		c.BodyBuffer = bytes.NewBufferString(callback + "(" + string(buf) + ")")
	default:
		c.Body = location
	}
	return
}
