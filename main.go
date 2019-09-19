package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vicanso/elton"
	_ "github.com/vicanso/location/controller"
	"github.com/vicanso/location/router"

	humanize "github.com/dustin/go-humanize"

	compress "github.com/vicanso/elton-compress"
	errorHandler "github.com/vicanso/elton-error-handler"
	etag "github.com/vicanso/elton-etag"
	fresh "github.com/vicanso/elton-fresh"
	recover "github.com/vicanso/elton-recover"
	responder "github.com/vicanso/elton-responder"
	stats "github.com/vicanso/elton-stats"
)

var (
	runMode string
)

// 获取监听地址
func getListen() string {
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":7001"
	}
	return listen
}

func main() {
	listen := getListen()

	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 只针对panic 以上的日志增加stack trace
	logger, err := c.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		panic(err)
	}

	d := elton.New()

	d.OnError(func(c *elton.Context, err error) {
		logger.DPanic("unexpected error",
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
	})

	d.Use(recover.New())

	d.Use(stats.New(stats.Config{
		OnStats: func(statsInfo *stats.Info, _ *elton.Context) {
			logger.Info("access log",
				zap.String("ip", statsInfo.IP),
				zap.String("method", statsInfo.Method),
				zap.String("uri", statsInfo.URI),
				zap.Int("status", statsInfo.Status),
				zap.String("consuming", statsInfo.Consuming.String()),
				zap.String("size", humanize.Bytes(uint64(statsInfo.Size))),
			)
		},
	}))

	d.Use(errorHandler.NewDefault())

	d.Use(func(c *elton.Context) error {
		c.NoCache()
		return c.Next()
	})

	d.Use(fresh.NewDefault())
	d.Use(etag.NewDefault())
	d.Use(compress.NewDefault())

	d.Use(responder.NewDefault())

	// health check
	d.GET("/ping", func(c *elton.Context) (err error) {
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
