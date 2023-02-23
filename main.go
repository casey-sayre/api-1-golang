package main

import (
	"context"
	"example/golang-api/config"
	"example/golang-api/controllers"
	pgRepo "example/golang-api/repositories/postgres"
  snsRepo "example/golang-api/repositories/sns"
	"example/golang-api/web"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {

	env := "local" // NOTE

  // log

	logger, _ := zap.NewDevelopment() // or NewProduction() or NewExample()
	slogger := logger.Sugar()

  // wire up the DI container

	fxApp := fx.New(
		fx.Provide(func() *config.Config {
			return config.NewConfig(env)
		}),
		fx.Provide(func(lc fx.Lifecycle) *zap.SugaredLogger {

			lc.Append(
				fx.Hook{
					OnStop: func(context.Context) error {
						slogger.Info("syncing logger")
						return slogger.Sync()
					},
				},
			)
			return slogger
		}),
		fx.Provide(pgRepo.ProvideAlbumRepository),
		fx.Provide(controllers.NewAlbumController),
		fx.Provide(web.NewFxRouter),
    fx.Provide(snsRepo.NewAlbumUpdatesPublisher),
		fx.Invoke(web.RegisterAlbumController),
	)

	// start the DI container

	go func() {
		startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := fxApp.Start(startCtx); err != nil {
			slogger.Fatal(err)
		}
	}()

	// listen for and wait for an interrupt signal

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	// stop the DI container

	if err := fxApp.Stop(context.Background()); err != nil {
		slogger.Error("fx stop, not graceful", err)
	} else {
		slogger.Info("fx stopping gracefully")
	}
}
