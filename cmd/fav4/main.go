package main

import (
	"context"
	"github.com/Snakdy/fav4/internal/api"
	"github.com/djcass44/go-utils/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
	"gitlab.com/autokubeops/serverless"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type environment struct {
	Port     int `envconfig:"PORT" default:"8080"`
	LogLevel int `split_words:"true"`
}

func main() {
	var e environment
	envconfig.MustProcess("app", &e)

	// configure logging
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(e.LogLevel * -1))

	log, _ := logging.NewZap(context.Background(), zc)

	// start the server
	route := api.NewIconAPI()
	serverless.NewBuilder(cors.AllowAll().Handler(route)).
		WithLogger(log).
		WithPort(e.Port).
		Run()
}
