package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vicanso/cod"
	"github.com/vicanso/cod/middleware"
	_ "github.com/vicanso/location/controller"
	"github.com/vicanso/location/router"
)

func main() {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":7001"
	}

	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 只针对panic 以上的日志增加stack trace
	logger, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}

	d := cod.New()

	d.Use(middleware.NewRecover())
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

	d.Use(middleware.NewResponder(middleware.ResponderConfig{}))

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
