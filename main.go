package main

import (
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/vicanso/elton"
	"github.com/vicanso/elton/middleware"
	_ "github.com/vicanso/location/controller"
	"github.com/vicanso/location/router"

	humanize "github.com/dustin/go-humanize"
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

	e := elton.New()

	e.OnError(func(c *elton.Context, err error) {
		logger.DPanic("unexpected error",
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
	})

	e.Use(middleware.NewRecover())

	e.Use(middleware.NewStats(middleware.StatsConfig{
		OnStats: func(statsInfo *middleware.StatsInfo, _ *elton.Context) {
			logger.Info("access log",
				zap.String("ip", statsInfo.IP),
				zap.String("method", statsInfo.Method),
				zap.String("uri", statsInfo.URI),
				zap.Int("status", statsInfo.Status),
				zap.String("latency", statsInfo.Latency.String()),
				zap.String("size", humanize.Bytes(uint64(statsInfo.Size))),
			)
		},
	}))

	e.Use(middleware.NewDefaultError())

	e.Use(func(c *elton.Context) error {
		c.NoCache()
		return c.Next()
	})

	e.Use(middleware.NewDefaultFresh())
	e.Use(middleware.NewDefaultETag())

	e.Use(middleware.NewDefaultResponder())

	// health check
	e.GET("/ping", func(c *elton.Context) (err error) {
		c.Body = "pong"
		return
	})

	groups := router.GetGroups()
	for _, g := range groups {
		e.AddGroup(g)
	}

	// http1与http2均支持
	e.Server = &http.Server{
		Handler: h2c.NewHandler(e, &http2.Server{}),
	}

	logger.Info("server will listen on " + listen)
	err = e.ListenAndServe(listen)
	if err != nil {
		panic(err)
	}
}
