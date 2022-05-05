package server

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infoidx/logger"
	"github.com/infoidx/server/middleware"
	"go.uber.org/fx"
)

type Option func(engine *gin.Engine) error

var UseRecovery = func() Option {
	return func(engine *gin.Engine) error {
		if engine == nil {
			return ErrGinInstanceNotInit
		}
		engine.Use(gin.Recovery())
		return nil
	}
}

var UseLogger = func(writer io.Writer) Option {
	return func(engine *gin.Engine) error {
		if engine == nil {
			return ErrGinInstanceNotInit
		}
		engine.Use(gin.LoggerWithWriter(writer))
		return nil
	}
}

var UseLogrus = UseLogger(logger.WithContext(context.Background()).Writer())

var UseCustomLogger = func() Option {
	return func(engine *gin.Engine) error {
		if engine == nil {
			return ErrGinInstanceNotInit
		}
		engine.Use(middleware.CustomLogger())
		return nil
	}
}

var UseCors = func() Option {
	return func(engine *gin.Engine) error {
		if engine == nil {
			return ErrGinInstanceNotInit
		}
		engine.Use(middleware.Cors())
		return nil
	}
}

var SetMode = gin.SetMode

func NewGinServer(opts ...Option) *gin.Engine {
	server := gin.New()
	for _, opt := range opts {
		err := opt(server)
		if err != nil {
			panic(err)
		}
	}
	// 增加两个默认路由
	server.GET("/ready", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	server.GET("/healthy", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	return server
}

func DefaultGinServer() *gin.Engine {
	return NewGinServer(UseLogrus, UseRecovery(), UseCors())
}

// 代码挪至具体业务中实现
func newGinFxLifeCycle(lc fx.Lifecycle, g *gin.Engine) {

}
