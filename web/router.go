package web

import (
	"context"
	"time"

	"example/golang-api/controllers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Router struct {
	Engine *gin.Engine // exported for testing
	server *http.Server
  slog *zap.SugaredLogger
}

func RegisterAlbumController(router *Router, albumController *controllers.AlbumController) {
	router.Engine.GET("/albums", albumController.GetAlbums)
	router.Engine.PATCH("/albums/:albumId", albumController.PatchAlbum)
}

func NewRouter(slogger *zap.SugaredLogger) *Router {
	router := Router{slog: slogger}
	router.Engine = gin.Default()
	router.server = &http.Server{
		Addr:    ":8080",
		Handler: router.Engine,
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Engine.Use(cors.New(corsConfig))

	router.Engine.SetTrustedProxies(nil)

	return &router
}

func NewFxRouter(lc fx.Lifecycle, slogger *zap.SugaredLogger) *Router {
  router := NewRouter(slogger)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {

			slogger.Info("Router fx OnStart: Starting http server...")

			// Initializing the server in a goroutine so that
			// it won't block the graceful shutdown handling below
			go func() {
				if err := router.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slogger.Fatalf("Router ListenAndServe failed: %s\n", err)
				}
				slogger.Info("Router ListenAndServe has shut down")
			}()

			return nil
		},
		OnStop: func(context.Context) error {

			slogger.Info("Router fx OnStop: Shutting down http server...")

			// The context is used to inform the server it has 5 seconds to finish
			// the request it is currently handling
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := router.server.Shutdown(ctx); err != nil {
				slogger.Fatal("Router http Server forced shutdown: ", err)
			}

			slogger.Info("Router http Server exiting")
			return nil
		},
	})
	return router
}
