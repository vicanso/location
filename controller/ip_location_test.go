package controller

import (
	"net/http/httptest"
	"testing"

	"github.com/vicanso/cod"
	"github.com/vicanso/location/service"
)

func TestGetLocation(t *testing.T) {
	ctrl := ipLocationCtrl{}
	t.Run("get client location", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ip-location/json/127.0.0.1", nil)
		req.Header.Set(cod.HeaderXForwardedFor, "1.2.7.255")
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		c.Params = map[string]string{
			"ip": "127.0.0.1",
		}
		err := ctrl.getLocation(c)
		if err != nil {
			t.Fatalf("get client location fail, %v", err)
		}
		location := c.Body.(*service.IPLocation)
		if location.City != "福州市" {
			t.Fatalf("get location fail")
		}
	})

	t.Run("get unknown ip", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ip-location/json/192.168.1.1", nil)
		resp := httptest.NewRecorder()
		c := cod.NewContext(resp, req)
		c.Params = map[string]string{
			"ip": "192.168.1.1",
		}
		err := ctrl.getLocation(c)
		if err != nil {
			t.Fatalf("get client location fail, %v", err)
		}
		location := c.Body.(*service.IPLocation)
		if location != unknownIPLocation {
			t.Fatalf("get unknown location fail")
		}
	})
}
