package controller

import (
	"os"

	"github.com/gobuffalo/packr"
	"github.com/vicanso/cod"
	"github.com/vicanso/cod/middleware"
	"github.com/vicanso/location/router"
)

type (
	// assetCtrl asset ctrl
	assetCtrl struct {
	}
	staticFile struct {
		box packr.Box
	}
)

var (
	box = packr.NewBox("../web/build")
)

func (sf *staticFile) Exists(file string) bool {
	return sf.box.Has(file)
}
func (sf *staticFile) Get(file string) ([]byte, error) {
	return sf.box.Find(file)
}
func (sf *staticFile) Stat(file string) os.FileInfo {
	return nil
}

func init() {
	g := router.NewGroup("")
	ctrl := assetCtrl{}
	g.GET("/", ctrl.index)

	sf := &staticFile{
		box: box,
	}
	g.GET("/static/*file", middleware.NewStaticServe(sf, middleware.StaticServeConfig{}))
}

func (ctrl assetCtrl) index(c *cod.Context) (err error) {
	file := "index.html"
	html, err := box.Find(file)
	if err != nil {
		return
	}
	c.SetFileContentType(file)
	c.Body = html
	return
}
