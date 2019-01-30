package main

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vicanso/cod"
	"github.com/vicanso/cod/middleware"
	_ "github.com/vicanso/location/controller"
	"github.com/vicanso/location/router"
)

// 获取监听地址
func getListen() string {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":7001"
	}
	return listen
}

func check() {
	listen := getListen()
	url := ""
	if listen[0] == ':' {
		url = "http://127.0.0.1" + listen + "/ping"
	} else {
		url = "http://" + listen + "/ping"
	}
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
		return
	}
	os.Exit(0)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "check" {
		check()
		return
	}
	listen := getListen()

	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 只针对panic 以上的日志增加stack trace
	logger, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}

	d := cod.New()

	d.Use(middleware.NewRecover())
	d.Use(middleware.NewDefaultErrorHandler())

	d.Use(middleware.NewStats(middleware.StatsConfig{
		OnStats: func(statsInfo *middleware.StatsInfo, _ *cod.Context) {
			logger.Info("access log",
				zap.String("ip", statsInfo.IP),
				zap.String("method", statsInfo.Method),
				zap.String("uri", statsInfo.URI),
				zap.Int("status", statsInfo.Status),
				zap.String("consuming", statsInfo.Consuming.String()),
			)
		},
	}))

	d.Use(middleware.NewDefaultCompress())

	d.Use(middleware.NewDefaultResponder())

	// health check
	d.GET("/ping", func(c *cod.Context) (err error) {
		c.Body = "pong"
		return
	})

	groups := router.GetGroups()
	for _, g := range groups {
		d.AddGroup(g)
	}

	logger.Info("server will listen on " + listen)
	err = d.ListenAndServe(listen)
	if err != nil {
		panic(err)
	}
}
