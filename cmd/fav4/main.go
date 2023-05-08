package main

import (
	"context"
	"github.com/djcass44/go-tracer/tracer"
	"github.com/djcass44/go-utils/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/cors"
	"gitlab.com/autokubeops/serverless"
	"gitlab.dcas.dev/open-source/fav4/internal/api"
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

	log, _ := logging.NewZap(context.TODO(), zc)

	// start the server
	route := api.NewIconAPI()
	serverless.NewBuilder(tracer.NewHandler(cors.AllowAll().Handler(route))).
		WithLogger(log).
		WithPort(e.Port).
		Run()
}
