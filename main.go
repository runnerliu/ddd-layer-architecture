package main

import (
	"context"
	"ddd-demo/infrastructure/config"
	"ddd-demo/interface/web/gin/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
)

var ginRouter *router.GinRouter

// PreStart 启动 HTTP Server 前的准备工作
func PreStart(conf config.Configuration) {
	ginRouter = router.NewGinRouter(
		conf.GetString("server.appName"),
		conf.GetBool("server.logWithBody"),
	)
	ginRouter.Start()
}

// ServeHTTP 启动以及关闭 HTTP Server
func ServeHTTP(lc fx.Lifecycle, conf config.Configuration) {
	serverPort := conf.GetInt("server.port")
	if serverPort <= 0 {
		serverPort = 8080
	}
	serverShutdownTimeoutMS := conf.GetInt("server.shutdownTimeoutTS")
	if serverShutdownTimeoutMS <= 0 {
		serverShutdownTimeoutMS = 1000
	}
	var srv *http.Server
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			srv = &http.Server{
				Addr:    fmt.Sprintf(":%d", serverPort),
				Handler: router.Router,
			}

			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			var timeout = time.Duration(serverShutdownTimeoutMS)
			ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Millisecond)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				panic(err)
			}

			select {
			case <-ctx.Done():
				fmt.Printf("%d milliseconds reached", timeout)
			}

			fmt.Println("Server exited")
			return nil
		},
	})
}

// @title DDD 分层架构 Demo
// @version 1.0
// @description DDD 分层架构 Demo
// @contact.name DDD-Layer-Architecture
// @contact.url http://www.swagger.io/support
// @contact.email william@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host http://localhost:8080
// @BasePath /api
// @Schemes http https
func main() {
	app := fx.New(
		fx.Provide(
			config.NewYamlConfiguration,
		),
		fx.Invoke(
			PreStart,
			ServeHTTP,
		),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := app.Stop(context.Background()); err != nil {
		panic(err)
	}
}
