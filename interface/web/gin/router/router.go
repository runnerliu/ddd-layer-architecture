package router

import (
	"ddd-demo/interface/web/gin/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	appName     string
	logWithBody bool
}

var (
	Router *gin.Engine
	ApiV1  *gin.RouterGroup
)

func NewGinRouter(appName string, logWithBody bool) *GinRouter {
	return &GinRouter{
		appName:     appName,
		logWithBody: logWithBody,
	}
}

func (r *GinRouter) Start() *gin.Engine {
	Router = gin.New()
	// 添加 CORS 中间件
	Router.Use(middleware.CORS())
	// 添加 GinLog
	Router.Use(middleware.GinLog(r.logWithBody))
	// 添加健康检查接口
	Router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	ApiV1 = Router.Group("/api/v1")

	return Router
}
